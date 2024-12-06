package main

import (
	"net/http"
)


func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{
        "status": "available",
        "environment": app.config.env,
        "version": version,
    }

    // 使用自定义模块helper中的writeJSON来转换数据
    err := app.writeJSON(w, http.StatusOK, data, nil)

    if err != nil {
        app.logger.Print(err)
        http.Error(w, "The server encoutered a problem and could not process your request", http.StatusInternalServerError)
    }
}
