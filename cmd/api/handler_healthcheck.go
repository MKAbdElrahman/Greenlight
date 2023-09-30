package main

import (
	"net/http"

	"greenlight.mkabdelrahman.net/internal/jsonparser"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	if err := jsonparser.WriteJSON(w, http.StatusOK, env, nil); err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
		return
	}
}
