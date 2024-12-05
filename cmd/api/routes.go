package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)


func (app *application) routes() *httprouter.Router {
    // 初始化httprouter实例
    router := httprouter.New()

    router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
    router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
    router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

    return router
    
}