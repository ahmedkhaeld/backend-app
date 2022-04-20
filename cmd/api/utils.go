package main

import (
	json2 "encoding/json"
	"net/http"
)

//writeJSON write json to the browser, takes in w , r, status code, the data to convert to json, and
// wrap which is string to wrap the json with a key that describe the content that comes out
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {

	wrapper := make(map[string]interface{})

	wrapper[wrap] = data

	json, err := json2.Marshal(wrapper)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(json)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}
	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{
		Message: err.Error(),
	}

	app.writeJSON(w, statusCode, theError, "error")
}
