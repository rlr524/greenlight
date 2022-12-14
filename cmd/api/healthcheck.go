package main

import (
	"fmt"
	"net/http"
	"os"
)

// TODO: return fqdn as well; currently returns host name
//func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
//	// A fixed-format JSON response from a string. We backticks in order to enclose double-quote characters
//	// in the JSON without needing to escape them, and the %q verb from the fmt package, which designates
//	// a double-quoted string that is safely escaped.
//	host, err := os.Hostname()
//	if err != nil {
//		panic(err)
//	}
//	js := `{"status":"available", "fqdn":%q "environment":%q, "version":%q}`
//	js = fmt.Sprintf(js, host, app.config.env, version)
//
//	// Set the content-type header on the response to application-json
//	w.Header().Set("Content-Type", "application/json")
//
//	_, e := w.Write([]byte(js))
//	if e != nil {
//		http.NotFound(w, r)
//		return
//	}
//}

// Refactoring the healthcheck handler to use a map of the healthcheck data and the writeJSON helper
// function which uses the json.Marshal function to return Go native objects as JSON text
func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// Get the host name
	host, e := os.Hostname()
	if e != nil {
		fmt.Println("the host name could not be resolved")
	}

	// A map that holds the information to return to the healthcheck, contained in an envelope as
	// defined in the writeJSON helper.
	envelope := envelope{
		"status": "available",
		"system_info": map[string]string{
			"fqdn":        host,
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, envelope, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
