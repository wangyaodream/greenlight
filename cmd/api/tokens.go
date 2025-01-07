package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/wangyaodream/greenlight/internal/data"
	"github.com/wangyaodream/greenlight/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	// 从请求中解析email和password
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// 验证email和password字段
	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// 通过email检索用户
	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// 验证密码
	match, err := user.Password.Matchs(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// 如果密码不匹配，则返回一个401 Unauthorized响应
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	// 创建一个新的authentication token
	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// 将token实例编码为JSON
	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
