package main
import (
	"strconv"
	"bytes"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/container"
	"crypto/rand"
	"io/ioutil"
	"os"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/hyperledger/fabric/core/crypto/primitives/ecies"
	"github.com/hyperledger/fabric/core/util"
	pb "github.com/hyperledger/fabric/protos"
	"tool/loadKey"
	"tool/rpc"
	"tool/initViper"
	"tool/transaction"
)

const (
	localStore string = "/var/hyperledger/production/client/"
)

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


/*

func getChaincodeBytes(context context.Context, spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	var codePackageBytes []byte
	var err error
	codePackageBytes, err = container.GetChaincodePackageBytes(spec)
	if err != nil {
		return nil, err
	}
	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes}
	return chaincodeDeploymentSpec, nil
}

func getMetadata(chaincodeSpec *pb.ChaincodeSpec) ([]byte, error) {
	return chaincodeSpec.Metadata, nil
}


func encryptTx(tx *pb.Transaction) error {
	eciesSPI := ecies.NewSPI()
	ccPrivateKey, err := eciesSPI.NewPrivateKey(rand.Reader, primitives.GetDefaultCurve())
	if err != nil {
		panic(fmt.Errorf("Failed generete chaincode keypair: %v\n", err))
		return err
	}

	var (
		stateKey   []byte
		privBytes []byte
	)

	switch tx.Type {
		case pb.Transaction_CHAINCODE_DEPLOY:
		  stateKey, err = primitives.GenAESKey()
		  if err != nil {
			panic(fmt.Errorf("Failed creating state key: %v\n", err))
			return err
		  }

		  privBytes, err = eciesSPI.SerializePrivateKey(ccPrivateKey)
		  if err != nil {
			panic(fmt.Errorf("Failed serializing chaincode key: %v\n", err))
			return err
		  }
		  break

		case pb.Transaction_CHAINCODE_INVOKE:
		  stateKey   = make([]byte, 0)
		  privBytes, err = eciesSPI.SerializePrivateKey(ccPrivateKey)
		  if err != nil {
			return err
		  }
		  break
	}

	chainPublicKey, err := loadKey.LoadKey()
	if err != nil {
		fmt.Println("error")
	}

	//fmt.Println("This is chainPublicKey: ", chainPublicKey.pub)

	cipher, err := eciesSPI.NewAsymmetricCipherFromPublicKey(chainPublicKey)

	msgToValidators, err := asn1.Marshal(chainCodeValidatorMessage1_2{privBytes, stateKey})
	if err != nil {
		panic(fmt.Errorf("Failed to preparing message to the validators: %v", err))
	}

	encMsgToValidators, err := cipher.Process(msgToValidators)
	if err != nil {
		panic(fmt.Errorf("Failed to encrypting message to the validators: %v", err))
	}
	tx.ToValidators = encMsgToValidators

	//initilize a new cipher
	cipher, err = eciesSPI.NewAsymmetricCipherFromPublicKey(ccPrivateKey.GetPublicKey())
	if err != nil {
		panic(fmt.Errorf("Failed initliazing encryption scheme: %v", err))
	}

	encryptedChaincodeID, err := cipher.Process(tx.ChaincodeID)
	if err != nil {
		panic(fmt.Errorf("Failed encrypting chaincodeID: %v", err))
	}
	tx.ChaincodeID = encryptedChaincodeID

	encryptedPayload, err := cipher.Process(tx.Payload)
	if err != nil {
		panic(fmt.Errorf("Failed encrypting payload: %v", err))
	}
	tx.Payload = encryptedPayload

	if len(tx.Metadata) != 0 {
		encryptedMetadata, err := cipher.Process(tx.Metadata)
		if err != nil {
			panic(fmt.Errorf("Failed to encrypt metadata"))
		}
		tx.Metadata = encryptedMetadata
	}

	return nil
}
*/
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
	for i := 0; i < n; i++ {
		tx := MakeInvokeTx()
		transactions = append(transactions, tx)
	}
		query := MakeQueryTx()
		//transactions = append(transactions, query)
//add a query transaction after finished ctrating invoke transactions
	for _, tx := range transactions {
		go func () {
		   _ = rpc.Connect(tx)
		   c <-1
		}()
		//fmt.Println(tx.Nonce)
	}

	for s := 0 ; s < n ; {
		s += <-c
	}

	response := rpc.Connect(query)
	fmt.Println(response)
}



