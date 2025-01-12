package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/wangyaodream/greenlight/internal/data"
	"github.com/wangyaodream/greenlight/internal/validator"
	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// 利用一个后台的 goroutine 定期清理过期的客户端
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				// 如果客户端在过去的 3 分钟内没有被看到过，那么就将它从map中删除
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		// 记录客户端的最后访问时间
		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求头中提取 Authorization 头
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]
		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			// 如果token不是26个字符长，则返回一个400 Bad Request响应
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// 通过token检索用户
		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				// 如果token无效，则返回一个401 Unauthorized响应
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				// 如果在检索用户时发生其他错误，则返回一个500 Internal Server Error响应
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

// func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
// 	// 这个中间件函数接受一个 http.HandlerFunc 作为参数，然后返回一个新的 http.HandlerFunc
// 	// 这样可以包装/v1/movie**路由处理函数
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		user := app.contextGetUser(r)

// 		if user.IsAnonymous() {
// 			app.authenticationRequireResponse(w, r)
// 			return
// 		}

// 		if !user.Activated {
// 			app.inactiveAccountResponse(w, r)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// 将requireActivatedUser中间件分离成验证用户和激活用户两个中间件
func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			// 如果用户未通过身份验证，则返回一个 401 Unauthorized 响应
			app.authenticationRequireResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if !user.Activated {
			// 如果用户未激活，则返回一个 403 Forbidden 响应
			app.inactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})

	return app.requireAuthenticatedUser(fn)
}

// 许可中间件
func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
    fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        user := app.contextGetUser(r)

        permissions, err := app.models.Permissions.GetAllForUser(user.ID)
        if err != nil {
            app.serverErrorResponse(w, r, err)
            return
        }

        if !permissions.Include(code) {
            app.notPermitedResponse(w, r)
            return
        }

        next.ServeHTTP(w, r)

    })

    return app.requireActivatedUser(fn)
}
