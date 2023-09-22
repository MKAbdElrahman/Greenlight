package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	if err := app.writeJSON(w, http.StatusOK, env, nil); err != nil {
		app.logger.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
