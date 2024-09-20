package main

import (
	internalHttp "diceroller_server/http"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	urlRouter := mux.NewRouter()

	// urlRouter.HandleFunc("/", homeHandler)

	urlRouter.HandleFunc("/diceroller", internalHttp.RollDice).Methods(http.MethodGet)
	urlRouter.HandleFunc("/token", internalHttp.HandleTokenGeneration).Methods(http.MethodPost, http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", urlRouter))
}
