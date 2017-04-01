package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	//"math/rand"
	//"time"
)

type response struct{
	GpCoin string      `json:"gpcoin,omitempty"`
	USD    string	    `json:"usd,omitempty"`
}

type User struct{
	GpCoin int
	USD    int
}




func writeHead(w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //
	w.Header().Set("content-type", "application/json")
	return w
	}

// login handler
func queryHandler(w http.ResponseWriter, request *http.Request) {
	w=writeHead(w)
	request.ParseForm()
	user, found := request.Form["User"]
	if !found {
		fmt.Fprint(w, "not find any user!")
		return
	}

	res := response{
		GpCoin : USERs[user[0]].GpCoin,
		USD : USERs[user[0]].USD,
		}
	b, err := json.Marshal(&res)
	fmt.Println(b)
	if err == nil{
		fmt.Fprint(w, string(b))
	}else{
		fmt.Fprint(w, err)
		return
	}
}

func topupHandler(w http.ResponseWriter, request *http.Request) {
	w=writeHead(w)
	request.ParseForm()
	user, found1  := request.Form["User"]
	amount, found2:= request.Form["Amount"]
	if !(found1 && found2) {
		fmt.Fprint(w, "not find any user!")
		return
	}
	usd, _ := strconv.Atoi(amount[0])
	var tmp =  User{
		GpCoin: USERs[user[0]].GpCoin,
		USD:    USERs[user[0]].USD + usd ,
	}

	USERs[user[0]] = tmp
}

/*
func investHandler(w http.ResponseWriter, request *http.Request) {
	writeHead(w)
	request.ParseForm()
	
	
}

func cashoutHandler(w http.ResponseWriter, request *http.Request) {
	writeHead(w)
	request.ParseForm()
	
	
}
*/

func transferHandler(w http.ResponseWriter, request *http.Request) {
	writeHead(w)
	request.ParseForm()
	from, found1  := request.Form["From"]
	to, found2    := request.Form["To"]
	amount, found3:= request.Form["Amount"]

	if !(found1 && found2 && found3) {
		fmt.Fprint(w, "not find any user!")
		return
	}

	trans, _ := strconv.Atoi(amount[0])

	if trans > USERs[from[0]].GpCoin{
		fmt.Fprint(w, "You don't have enough money!")
		return
		}
	var fromTmp =  User{
		GpCoin: USERs[from[0]].GpCoin - trans,
		USD:    USERs[from[0]].USD ,
	}

  var toTmp =  User{
		GpCoin: USERs[to[0]].GpCoin + trans,
		USD:    USERs[to[0]].USD ,
	}
	USERs[from[0]] = fromTmp
	USERs[to[0]] = toTmp

}

//home page handler

// server main method

var router = mux.NewRouter()

func main() {

	router.HandleFunc("/query", queryHandler).Methods("POST")
	router.HandleFunc("/topup", topupHandler).Methods("POST")
	//router.HandleFunc("/invest", investHandler).Method("POST")
	//router.HandleFunc("/cashout", cashoutHandler).Method("POST")
	router.HandleFunc("/transfer", transferHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":8000", nil)
}
