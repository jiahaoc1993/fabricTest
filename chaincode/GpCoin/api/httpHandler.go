package main
import(
//	"encodeing/json"
//	"io"
	"fmt"
	"net/http"
	"os"
)

var chaincodeName string


func login(w http.ResponseWriter, req *http.Request) {
	//fmt.Println()
	//req.ParseForm()
	//param_user, found1 := req.Form["user"]
	//param_pass, found2 := req.Form["pass"]

	//if !(found1 && found2) {
	//	fmt.Fprintf(w, "Don't do this")
	//
	req.ParseForm()
	userNames, found1 := req.Form["userName"]
	passs, found2     := req.Form["pass"]
	userName := userNames[0]
	pass    := passs[0]

	if !(found1 && found2) {
		fmt.Fprintf(w, "Please offer the user and pass")
	}

	if userName == "Tom" && pass == "123" {
		fmt.Fprintf(w, "ok")
	}else if userName == "alice" && pass == "123" {
		fmt.Fprintf(w, "ok")
	}else if userName == "bob"   &&  pass == "123" {
		fmt.Fprintf(w, "ok")
	}else {
		fmt.Fprintf(w, "user or pass error!")
	}

}

func topupHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	amount, found1 := req.Form["amount"]
	user, found2   := req.Form["user"]

	if !(found1 && found2) {
		fmt.Fprintf(w, "Wrong arguments")
	}

	err := topup(chaincodeName, amount[0], user[0])
	if err != nil {
		fmt.Fprintf(w, "faild")
	}

	fmt.Fprintf(w, "ok")
}


func investHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	amount, found1 := req.Form["amount"]
	user, found2   := req.Form["user"]

	if !(found1 && found2) {
		fmt.Fprintf(w, "Wrong arguments")
	}

	err := invest(chaincodeName, amount[0], user[0])
	if err != nil {
		fmt.Fprintf(w, "faild")
	}

	fmt.Fprintf(w, "ok")
}


func cashoutHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	amount, found1 := req.Form["amount"]
	user, found2   := req.Form["user"]

	if !(found1 && found2) {
		fmt.Fprintf(w, "Wrong arguments")
	}

	err := cashout(chaincodeName, amount[0], user[0])
	if err != nil {
		fmt.Fprintf(w, "faild")
	}

	fmt.Fprintf(w, "ok")
}

func transferHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	amount, found1 := req.Form["amount"]
	from, found2   := req.Form["from"]
	to, found3     := req.Form["to"]

	if !(found1 && found2 && found3) {
		fmt.Fprintf(w, "Wrong arguments")
	}

	err := transfer(chaincodeName, amount[0], from[0], to[0])
	if err != nil {
		fmt.Fprintf(w, "faild")
	}

	fmt.Fprintf(w, "ok")
}

func queryHandle(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	user, found := req.Form["user"]
	if !found {
		fmt.Fprintf(w, "wrong arguments")
	}

	res, err := CheckUser(chaincodeName, user[0])
	if err != nil {
		fmt.Fprintf(w, "failed")
	}

	fmt.Fprintf(w, res)
}

func main() {
	var err error
	if err = initNVP(); err != nil {
	appLogger.Debugf("Failed initiliazing NVP [%s]", err)
                os.Exit(-1)
	}

	chaincodeName, err = deploy()
	if err != nil {
		appLogger.Debugf("Failed with initiliazing")
		os.Exit(-1)
	}

	http.HandleFunc("/login", login)
	http.HandleFunc("/topup", topupHandle)
	http.HandleFunc("/invest", investHandle)
	http.HandleFunc("/cashout", cashoutHandle)
	http.HandleFunc("/transfer", transferHandle)
	http.HandleFunc("/query", queryHandle)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}


