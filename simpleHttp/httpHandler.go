package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// cookie handling



// login handler
func homePageHandler(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Hello World"))
}



//home page handler

// server main method

var router = mux.NewRouter()

func main() {

	router.HandleFunc("/query", homePageHandler)

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}