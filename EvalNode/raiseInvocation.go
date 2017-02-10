package main
import (
	"strconv"
	"encoding/json"
	"fmt"
	"os"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/hyperledger/fabric/core/util"
	pb "github.com/hyperledger/fabric/protos"
	"tool/loadKey"
	"tool/rpc"
	"tool/initViper"
	"tool/transaction"
	"time"
)

const (
	localStore string = "/var/hyperledger/production/client/"
)

type response struct {
	Name   string    `json:"name,omitempty"`
	Amount string    `json:"amount,omitemty"`
	Time   string    `json:"time,omitempty"`
}


type chainCodeValidatorMessage1_2 struct {
	PrivateKey []byte
	StateKey   []byte
}


func Init() (err error) { //init the crypto layer
	securityLevel := 256
	hashAlgorithm := "SHA3"
	if err = primitives.InitSecurityLevel(hashAlgorithm, securityLevel); err != nil {
		panic(fmt.Errorf("Failed setting security level: %v", err))
		return err
	}

	return nil
}


func Sign(tx *pb.Transaction) (*pb.Transaction, error) {
	enrollmentCert, privKey, err := loadKey.LoadEnrollment()
	if err != nil {
		fmt.Printf("Failed loading enrollment metieral")
		return nil, err
	}

	tx.Cert = enrollmentCert.Raw

	rawTx, err := proto.Marshal(tx)
	if err != nil {
		fmt.Printf("Failed marshaling tx: %v", err)
		return nil, err
	}

	rawSignature, err := primitives.ECDSASign(privKey, rawTx)
	if err != nil {
		fmt.Println("Failed Creating signature: %v", err)
		return nil, err
	}

	tx.Signature = rawSignature

	return tx, nil
}

func MakeInvokeTx() *pb.Transaction {
	chaincodeInvocationSpec, err := transaction.InvokeChaincodeSpec()
	if err != nil {
		os.Exit(0)
	}
	//fmt.Println(chaincodeInvocationSpec)

	if err != nil {
		os.Exit(0)
	}

	tx, err := transaction.CreateInvokeTx(chaincodeInvocationSpec, util.GenerateUUID(), nil, chaincodeInvocationSpec.ChaincodeSpec.Attributes...)
	if err != nil {
		os.Exit(0)
	}

	tx, err = Sign(tx)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	return tx
}

func MakeQueryTx() *pb.Transaction {
	chaincodeInvocationSpec, err := transaction.QueryChaincodeSpec()
	if err != nil {
		os.Exit(0)
	}
	//fmt.Println(chaincodeInvocationSpec)

	if err != nil {
		os.Exit(0)
	}

	tx, err := transaction.CreateQueryTx(chaincodeInvocationSpec, util.GenerateUUID(), nil, chaincodeInvocationSpec.ChaincodeSpec.Attributes...)
	if err != nil {
		os.Exit(0)
	}

	tx, err = Sign(tx)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	return tx
}

func QueryBeforeInvoke() string {
	var log    string
	var result response
	query := MakeQueryTx()
	res := rpc.Connect(query)
	t := strconv.FormatInt(time.Now().UnixNano(), 10)
	err := json.Unmarshal(res.Msg, &result)
	if err != nil {
		panic(fmt.Errorf("Failed unmarshaling json: %V\n", err))
	}
	log = fmt.Sprintln("***************")
	log = log + fmt.Sprintf("###Start at time: %s, while account %s has amount %s\n", t, result.Name, result.Amount)
	return log
}

func QueryAfterInvoke() string{
	var log string
	time.Sleep(2 * time.Second)
	var result response
	query := MakeQueryTx()
	res := rpc.Connect(query)
	err := json.Unmarshal(res.Msg, &result)
	if err != nil {
		panic(fmt.Errorf("Failed unmarshaling json: %V\n", err))
	}

	log = fmt.Sprintln("###Account " + result.Name + " have amount " + result.Amount + " at time: "+result.Time)
	return log
}



func main() {
	Init()
	err := initViper.SetConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading viper config file"))
	}
	n, err := strconv.Atoi(os.Args[1])
	c := make(chan int)
	transactions := []*pb.Transaction{}
	if err != nil {
		panic(fmt.Errorf("Failed conversing args"))
		}

	f, err := os.OpenFile("./log/"+ os.Args[1] +"Ivocations.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()

		//query := MakeQueryTx()
		//transactions = append(transactions, query)
//add a query transaction after finished ctrating invoke transactions
	//query the current state first and set the time as well!
	if n != 0 && os.Args[2] == "start" {
		before := QueryBeforeInvoke()
		fmt.Print(before)
		_, err = f.WriteString(before)
		if err != nil {
		panic(err)
		}


		for i := 0; i < n; i++ {
			tx := MakeInvokeTx()
			transactions = append(transactions, tx)
		}

		time.Sleep(time.Second) // time out the batch
		for _, tx := range transactions {
			go func () {
			_ = rpc.Connect(tx)
			 c <-1
			}()
		}
		for s := 0 ; s < n ; {
			s += <-c
		}
		//wait for thr processTransaction
	}else if os.Args[2] == "query"{

		after := QueryAfterInvoke()
		fmt.Print(after)
		_, err = f.WriteString(after)
		if err != nil {
			panic(err)
		}
	}else{
	fmt.Println("Wrong Args! please!!")
	}
}



