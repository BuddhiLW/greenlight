package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	// fmt.Fprintln(w, "status: available")
	// fmt.Fprintf(w, "environmnent: %s\n", app.config.env)
	// fmt.Fprintf(w, "version: %s\n", version)

	// js := `{"status": "available", "environment": %q, "version": %q}`
	// js = fmt.Sprintf(js, app.config.env, version)
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	js, err := json.Marshal(data)

	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server could not comply with the request", http.StatusInternalServerError)
		return
	}

	// for easiness of reading the response
	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(js))
}
