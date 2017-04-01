package main
import(
	"encodeing/json"
//	"io"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"github.com/gorilla/mux"
)

var chaincodeName string

type response struct{
        GpCoin string      `json:"gpcoin,omitempty"`
        USD    string       `json:"usd,omitempty"`
}


func writeHead(w http.ResponseWriter) http.ResponseWriter {
        w.Header().Set("Access-Control-Allow-Origin", "*")             //
        w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //
        w.Header().Set("content-type", "application/json")
        return w
}

func topupHandle(w http.ResponseWriter, req *http.Request) {
	w = writeHead(w)
	req.ParseForm()
	amount, found1 := req.Form["Amount"]
	user, found2   := req.Form["User"]

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
	amount, found1 := req.Form["Amount"]
	user, found2   := req.Form["User"]

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
	amount, found1 := req.Form["Amount"]
	user, found2   := req.Form["User"]

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
	amount, found1 := req.Form["Amount"]
	from, found2   := req.Form["From"]
	to, found3     := req.Form["To"]

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
	user, found := req.Form["User"]
	if !found {
		fmt.Fprintf(w, "wrong arguments")
	}

	res, err := CheckUser(chaincodeName, user[0])
	if err != nil {
		fmt.Fprintf(w, "failed")
	}

	amounts := strings.Split(amounts, ",")

	r := response{
		GpCoin  : amounts[0]
		USD	: amounts[1]
	}

	b, err := json.Marshal(&r)
        fmt.Println(b)
        if err == nil{
                fmt.Fprint(w, string(b))
        }else{
                fmt.Fprint(w, err)
                return
	}
}

var router = mux.NewRouter()

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

	//http.HandleFunc("/login", login)
	router.HandleFunc("/topup", topupHandle).Methods("POST")
	router.HandleFunc("/invest", investHandle).Methods("POST")
	router.HandleFunc("/cashout", cashoutHandle).Methods("POST")
	router.HandleFunc("/transfer", transferHandle).Methods("POST")
	router.HandleFunc("/query", queryHandle).Methods("POST")
	router.Handle("/", router)


	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println(err)
	}
}


