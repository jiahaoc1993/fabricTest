package main
import (
	"fmt"
	"os"
	"github.com/hyperledger/fabric/core/peer"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/hyperledger/fabric/core/util"
	pb "github.com/hyperledger/fabric/protos"
	context "golang.org/x/net/context"
	"tool/loadKey"
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

func FakeSign(tx *pb.Transaction) (*pb.Transaction, error) {
	enrollmentCert, privKey, err := loadKey.LoadFakeEnrollment()
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

func MakeFakeTx(chaincodeName string, args []string) *pb.Transaction {
	chaincodeInvocationSpec, err := transaction.InvokeChaincodeSpec(chaincodeName, args)
	if err != nil {
		os.Exit(0)
	}
	//fmt.Println(chaincodeInvocationSpec)

	tx, err := transaction.CreateInvokeTx(chaincodeInvocationSpec, util.GenerateUUID(), nil, chaincodeInvocationSpec.ChaincodeSpec.Attributes...)
	if err != nil {
		os.Exit(0)
	}

	tx, err = FakeSign(tx)
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
//	var timeout           int
	flag.StringVar(&method, "m", "", "method of Execution")
	flag.StringVar(&chaincodeName, "n", "", "Name of chaincode returned by the deploy transaction")
	flag.IntVar(&numOfTransactions, "t", 1, "Number of transaction readly to send(dafault=1)")
//	flag.IntVar(&timeout, "timeout", 0, "fake transaction duration")
	Init()					//viper init
	err := initViper.SetConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading viper config file"))
	}
	//c := make(chan int)		       //main exit after all go rutines were lanuched
	transactions := []*pb.Transaction{}    //array of transactiongs

	flag.Parse()


	for i := 0; i < numOfTransactions; i++ {
		tx := MakeFakeTx(chaincodeName,[]string{"a","b","1"})
		transactions = append(transactions, tx)
	}

	con, err := peer.NewPeerClientConnectionWithAddress("10.0.2.15")
	if err != nil {
		fmt.Println("Can't not connect: %v", err)
	}

	defer con.Close()
	client := pb.NewPeerClient(con)


	for i := 0 ; i < numOfTransactions ; i++ {
		_, _ = client.ProcessTransaction(context.Background(), transactions[i])
	}
	fmt.Println("send 1000 spent time: ", spent )
}//main
