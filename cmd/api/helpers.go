package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/wangyaodream/greenlight/internal/validator"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	// Add the content type header.
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
    // 限制请求体的大小为 1MB
    maxBytes := 1_048_576
    r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

    dec := json.NewDecoder(r.Body)
    // 设置 decoder 为严格模式，不允许包含未知字段
    dec.DisallowUnknownFields()

    err := dec.Decode(dst)
    if err != nil {
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError
        var invalidUmarsalError *json.InvalidUnmarshalError
        var maxBytesError *http.MaxBytesError

        // errors.IS 函数检查错误是否是特定类型的错误
        // errors.As 函数检查错误是否实现了特定接口

        switch {
        case errors.As(err, &syntaxError):
            return fmt.Errorf("request body contains badly-formed JSON (at character %d)", syntaxError.Offset)
        case errors.Is(err, io.ErrUnexpectedEOF):
            return errors.New("request body contains badly-formed JSON")
        case errors.As(err, &unmarshalTypeError):
            if unmarshalTypeError.Field != "" { 
                return fmt.Errorf("request body contains an invalid value for the %q field", unmarshalTypeError.Field)
            }
            return fmt.Errorf("body contains an invalid value at character %d", unmarshalTypeError.Offset)
        case errors.Is(err, io.EOF):
            return errors.New("request body must not be empty")
        case errors.As(err, &maxBytesError):
            return fmt.Errorf("request body must not be larger than %d bytes", maxBytesError.Limit)
        case strings.HasPrefix(err.Error(), "json: unknown field "):
            fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
            return fmt.Errorf("request body contains unknown field %s", fieldName)

        case errors.As(err, &invalidUmarsalError):
            panic(err)


        default:
            return err

        }
    }
    // 检查请求体是否只包含一个 JSON 值
    err = dec.Decode(&struct{}{})
    if err != io.EOF {
        return errors.New("request body must only contain a single JSON value")
    }
    return nil
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
    // 从 URL 查询字符串中获取字符串值
    s := qs.Get(key)

    if s == "" {
        return defaultValue
    }

    return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
    // 从 URL 查询字符串中获取字符串值
    csv := qs.Get(key)

    if csv == "" {
        return defaultValue
    }

    return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
    s := qs.Get(key) 

    if s == "" {
        return defaultValue
    }

    i, err := strconv.Atoi(s)
    if err != nil {
        v.AddError(key, "must be an integer value")
        return defaultValue
    }

    return i
}

func (app *application) background(fn func()) {
    go func() {
        defer func() {
            if err := recover(); err != nil {
                // 在这里记录错误信息
                app.logger.PrintError(fmt.Errorf("%s", err), nil)
            }
        }()

        // 调用传入的函数
        fn()
    }()
}
