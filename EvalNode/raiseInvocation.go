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
	"flag"
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


func Init() (err error) { //init the crypto layer or use crypto.Init()
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

func MakeInvokeTx(chaincodeName string) *pb.Transaction {
	chaincodeInvocationSpec, err := transaction.InvokeChaincodeSpec(chaincodeName)
	if err != nil {
		os.Exit(0)
	}
	//fmt.Println(chaincodeInvocationSpec)

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

func MakeQueryTx(chaincodeName string, TT string) *pb.Transaction {
	chaincodeInvocationSpec, err := transaction.QueryChaincodeSpec(chaincodeName, TT)
	if err != nil {
		os.Exit(0)
	}
	//fmt.Println(chaincodeInvocationSpec)


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

func WarnningMsg() string{
//      var err error
        var str string
        str =  "Usage:\n     ./raiseInvocation [command]\n\n"
        str += "Available Commands:\n"
        str += "     invoke      Invoke the specified chaincode\n"
        str += "     query       Query the specified chaincode \ni\n"
        str += "Flags:\n     -n, --name string     Name of chaincode returned by the deploy teansaction"
        return str
}





func main() {
	var numOfTransactions int
	var chaincodeName     string
	var method	      string
	flag.StringVar(&method, "m", "", "method of Execution")
	flag.StringVar(&chaincodeName, "n", "", "Name of chaincode returned by the deploy transaction")
	flag.IntVar(&numOfTransactions, "t", 1, "Number of transaction readly to send(dafault=1)")
	Init()					//viper init
	err := initViper.SetConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading viper config file"))
	}
	c := make(chan int)		       //main exit after all go rutines were lanuched
	transactions := []*pb.Transaction{}    //array of transactiongs

	flag.Parse()

	if chaincodeName == "" || method == ""{    //chaincode name must need
		panic(fmt.Errorf("method of name of chaincode should not be empty"))
	}

	if method == "invoke" {
		var res response
		var stateBefore, stateAfter int
		var timeBefore, timeAfter float64

			//check the current state before taking invocation!
		query := MakeQueryTx(chaincodeName, "now")
		finish  := MakeQueryTx(chaincodeName, "state")
		response := rpc.Connect(query)
		_ = json.Unmarshal(response.Msg, &res)
		stateBefore, _ = strconv.Atoi(res.Amount)
		timeBefore , _ = strconv.ParseFloat(res.Time, 64)
		for i := 0; i < numOfTransactions; i++ {
			tx := MakeInvokeTx(chaincodeName)
			transactions = append(transactions, tx)
		}

		time.Sleep( 2 * time.Second) // time out the batch
		timeBefore = timeBefore + float64(2 * time.Second) // 
		for _, tx := range transactions {
			go func () {
			_ = rpc.Connect(tx)
			 c <-1
			}()
		}
		for s := 0 ; s < numOfTransactions ; {
			s += <-c
		}

		for i :=0; i < 10; i++ {
			time.Sleep(2 * time.Second)
			response = rpc.Connect(finish)
			 _ = json.Unmarshal(response.Msg, &res)
			stateAfter, _ = strconv.Atoi(res.Amount)
			if stateAfter == numOfTransactions + stateBefore {
					timeAfter, _ = strconv.ParseFloat(res.Time, 64)
					spent := (timeAfter - timeBefore) / float64(time.Second)
					fmt.Printf("Execute %d transactions spent %.3f seconds\n", numOfTransactions, spent)
					break
					}else if i < 9{
						continue
					}else {
						panic(fmt.Errorf("remote server run out of time to response!"))
					}
		}
	}else if method == "query"{
		query := MakeQueryTx(chaincodeName, "state")
		response := rpc.Connect(query)
		fmt.Println("Status: " + string(response.Status) + "," + "Msg: " + string(response.Msg))
	}else {
		panic(fmt.Errorf(WarnningMsg()))
	}
}//main
