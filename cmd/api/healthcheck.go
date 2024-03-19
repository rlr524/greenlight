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

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request",
			http.StatusInternalServerError)
	}

	//json := `{"status": "available", "environment": %q, "version": %q}`
	//json = fmt.Sprintf(json, app.config.env, version)
	//
	//w.Header().Set("Content-Type", "application/json")
	//_, err := w.Write([]byte(json))
	//if err != nil {
	//	app.logger.Error("There was an error writing the response: %v", err)
	//	return
	//}
}
