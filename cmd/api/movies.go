package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wangyaodream/greenlight/internal/data"
)


func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "create a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
    // 利用自定义的heaper模块，获取id参数
    id, err := app.readIDParam(r)

    if err != nil {
        app.notFoundResponse(w, r)
        return
    }

    movie := data.Movie{
        ID: id,
        Title: "Casablanca",
        CreatedAt: time.Now(),
        Year: 1942,
        Runtime: 102,
        Genres: []string{"Drama", "Romance", "War"},
        Version: 1,
    }

    err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err) 
    }

}
