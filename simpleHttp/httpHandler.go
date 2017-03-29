package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type response struct{
	GpCoin string       `json:"gpcoin,omitempty"`
	USD    string	    `json:"usd,omitempty"`
}

// cookie handling



// login handler
func homePageHandler(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //
	w.Header().Set("content-type", "application/json")
	res := response{
		GpCoin : "130",
		USD : "134",
		}
	b, err := json.Marshal(&res)
	fmt.Println(b)
	if err == nil{
		fmt.Fprint(w, string(b))
	}else{
		fmt.Fprint(w, err)
}	}



//home page handler

// server main method

var router = mux.NewRouter()

func main() {

	router.HandleFunc("/query", homePageHandler)

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
