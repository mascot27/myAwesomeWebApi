package main

import (
	"fmt"
	"net/http"

	ghandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mascot27/myAwesomeWebApi/handlers"
	"github.com/mascot27/myAwesomeWebApi/handlers/endpoints"
	"github.com/mascot27/myAwesomeWebApi/models"

	"log"
)

func main() {
	addressHttps := ":8443"
	certFileHttps := "server.crt"
	keyFileHttps := "server.key"

	fmt.Println("api accessible by: " + addressHttps)

	headersOk := ghandler.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := ghandler.AllowedOrigins([]string{"*"})
	methodsOk := ghandler.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// base server
	var router = mux.NewRouter().StrictSlash(true)

	// auth server
	var routerAuthServer = router.PathPrefix("/auth").Subrouter()
	routerAuthServer.HandleFunc("/get-token", authServerGetToken).Methods("POST")
	routerAuthServer.HandleFunc("/refresh-token", authServerRefreshToken).Methods("POST")


	// ROUTES - /api
	var routerApi = router.PathPrefix("/api").Subrouter()


	routerApi.HandleFunc("/home", homeHandler)

	// /api/get-token  (post with creds) and return token
	routerApi.HandleFunc("/get-token", endpoints.Authenticate).Methods("POST")

	routerApi.HandleFunc("/private", endpoints.Authenticate).Methods("GET")

	// server using https
	err := http.ListenAndServeTLS(addressHttps, certFileHttps, keyFileHttps,
		ghandler.CORS(originsOk, headersOk, methodsOk)(router))
	if err != nil {
		log.Fatal("listen and serve: ", err)
	}
}

func authServerGetToken(writer http.ResponseWriter, request *http.Request) {
	// 1) read user and password in request body
	// 2) validate credentials
	// 		- if bad -> bad request response
	//		- if good -> continue
	// 3)
	// 		- create a refresh token (long lived)
	//		- create an access token (short lived)
	// 		- save refresh token in whitelist
}

func authServerRefreshToken(writer http.ResponseWriter, request *http.Request) {
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Name:     "Corentin",
		Password: "password",
	}

	handlers.WriteJsonData(w, user, http.StatusOK)
}
