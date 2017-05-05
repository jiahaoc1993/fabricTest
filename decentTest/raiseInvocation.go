package main
import (
	"strconv"
	"encoding/json"
	"fmt"
	"os"
	"github.com/hyperledger/fabric/core/peer"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/hyperledger/fabric/core/util"
	pb "github.com/hyperledger/fabric/protos"
	context "golang.org/x/net/context"
	"tool/loadKey"
	//"tool/rpc"
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
func MakeInvokeTx(chaincodeName string, args []string) *pb.Transaction {
	chaincodeInvocationSpec, err := transaction.InvokeChaincodeSpec(chaincodeName, args)
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

func MakeQueryTx(chaincodeName string, args []string) *pb.Transaction {
	chaincodeInvocationSpec, err := transaction.QueryChaincodeSpec(chaincodeName, args)
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
	var dest	      string

	flag.StringVar(&dest, "d","","destination")
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

	if chaincodeName == ""{    //chaincode name must need
		panic(fmt.Errorf("name of chaincode should not be empty"))
	}
//start the invoke
	var res response
	var stateBefore, stateAfter int
	var timeBefore, timeAfter float64

	for i := 0; i < numOfTransactions; i++ {
		tx := MakeInvokeTx(chaincodeName,[]string{"a","b","1"})
		transactions = append(transactions, tx)
	}
			//check the current state before taking invocation!
	query := MakeQueryTx(chaincodeName, []string{"b"})
	finish  := MakeQueryTx(chaincodeName, []string{"b"})

	con, _ := peer.NewPeerClientConnectionWithAddress(dest)
	defer con.Close()
	client := pb.NewPeerClient(con)

	response,_ := client.ProcessTransaction(context.Background(), query)
	_ = json.Unmarshal(response.Msg, &res)
		stateBefore, _ = strconv.Atoi(res.Amount)
		timeBefore , _ = strconv.ParseFloat(res.Time, 64)

	//	time.Sleep( 2 * time.Second) // time out the batch
	//	timeBefore = timeBefore + float64(2 * time.Second) // 
	for _, tx := range transactions {
		go func () {
		_ , _= client.ProcessTransaction(context.Background(), tx)
			c <-1
		}()
	}

	for s := 0 ; s < numOfTransactions ; {
		s += <-c
	}

	for i :=0; i < 10; i++ {
		response,_ = client.ProcessTransaction(context.Background(), finish)
		_ = json.Unmarshal(response.Msg, &res)
		stateAfter, _ = strconv.Atoi(res.Amount)
		if stateAfter == numOfTransactions + stateBefore {
			timeAfter, _ = strconv.ParseFloat(res.Time, 64)
			spent := (timeAfter - timeBefore) / float64(time.Second)
			fmt.Printf("Execute %d transactions spent %.3f seconds\n", numOfTransactions, spent)
					return
		}
		time.Sleep(1 * time.Second)
	}

	panic(fmt.Errorf("remote server run out of time to response!"))
}//main
