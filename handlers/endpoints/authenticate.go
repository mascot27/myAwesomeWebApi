package endpoints

import (
	"encoding/json"
	middleware "myAwesomeWebApi/handlers/middlewares"
	"myAwesomeWebApi/models"
	"net/http"
)

func Authenticate(w http.ResponseWriter, r *http.Request) {

	var userRequest models.User
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&userRequest); err != nil {
		w.Write([]byte("Invalid request payload"))
		return
	}

	name := userRequest.Name
	password := userRequest.Password

	if len(name) == 0 || len(password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide name and password to obtain the token"))
		return
	}

	if models.IsValidCredentials(name, password) {
		token, err := middleware.GetToken(name)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error generating JWT token: " + err.Error()))
		} else {
			w.Header().Set("Authorization", "Bearer "+token)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Token: " + token))
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Name and password do not match"))
		return
	}
}
