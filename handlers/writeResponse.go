package handlers

import (
	"encoding/json"
	"net/http"
)

func WriteJsonData(writer http.ResponseWriter, data interface{}, status int) error {
	writer.Header().Set("Content-Type", "application/json; charset=UTF-8")
	writer.WriteHeader(status)
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		return err
	}
	return nil
}

/*
TODO: write standard errors
*/
