package main

import (
	"fmt"
	"net/http"

	ghandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"myAwesomeWebApi/handlers"
	"myAwesomeWebApi/handlers/endpoints"
	"myAwesomeWebApi/models"

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

	var router = mux.NewRouter().StrictSlash(true)

	var routerApi = router.PathPrefix("/api").Subrouter()

	routerApi.HandleFunc("/a", homeHandler)

	/*
		Authentication by JWT
	*/
	routerApi.HandleFunc("/get-token", endpoints.Authenticate).Methods("POST")

	/*
			User Ressource
		------------------------
			User registration
			User deletion
			User update
	*/

	/*
		resources's routes
	*/

	// server using https
	err := http.ListenAndServeTLS(addressHttps, certFileHttps, keyFileHttps,
		ghandler.CORS(originsOk, headersOk, methodsOk)(router))
	if err != nil {
		log.Fatal("listen and serve: ", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	user := models.User{
		Name:     "Corentin",
		Password: "password",
	}

	handlers.WriteJsonData(w, user, http.StatusOK)
}
