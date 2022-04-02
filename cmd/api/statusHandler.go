package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) statusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := AppStatus{
		Status:      "Available",
		Environment: app.config.env,
		Version:     version,
	}
	// convert to json
	js, err := json.MarshalIndent(currentStatus, "", "\t")
	if err != nil {
		app.logger.Println(err)
	}
	// write json to the browser in json format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
