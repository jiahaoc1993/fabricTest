package main
import (
	"strconv"
	"bytes"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/container"
	//"path/filepath"
	//	"github.com/hyperledger/fabric/core"
	//	"github.com/hyperledger/fabric/core/peer"
	//"github.com/spf13/viper"
	//"crypto/elliptic"
	"crypto/rand"
	"io/ioutil"
	"os"
	//"strings"
	//"errors"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	//"github.com/golang/protobuf/jsonpb"
	//"github.com/hyperledger/fabric/core/crypto"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/hyperledger/fabric/core/crypto/primitives/ecies"
	"github.com/hyperledger/fabric/core/util"
	pb "github.com/hyperledger/fabric/protos"
	"tool/loadKey"
	"tool/rpc"
	"tool/initViper"
)

const (
	info = "1600b975353e233708899a3b0ff8da55418f0738ef47f4a22d84b90da481d31261432209f1e7c4767dd2c400d685f4c96d41493e4576a52e41aaa36b142eaf81"
	localStore string = "/var/hyperledger/production/client/"
)

type chainCodeValidatorMessage1_2 struct {
	PrivateKey []byte
	StateKey   []byte
}

type Transaction struct {
	Jsonrpc string `json:"jsonrpc,omitempty"`
	Method  string `json:"method,omitempty"`
	Params  params `json:"params,omitempty"`
	ID      int    `json:"id,omitempty"`
}

type params struct {
	Type          int               `json:"type,omitempty"`
	ChaincodeID   map[string]string `json:"chaincodeID,omitempty"`
	CtorMsg       ctorMsg           `json:"ctorMsg"`
	SecureContext string            `josn:"secureContext,omitempty"`
}

type ctorMsg struct {
	Function string   `json:"function,omitempty"`
	Args     []string `json:"args,omitempty"`
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

// this is only for pb.chaincodSpec
func DeployChaincodeSpec() (*pb.ChaincodeSpec, error) {
	var spec pb.ChaincodeSpec
	//var spec2 pb.ChaincodeSpec
	t := &params{
		1,
		map[string]string{"path": "github.com/hyperledger/fabric/examples/chaincode/go/HelloWorld"},
		ctorMsg{"init", []string{"a"}},
		"lukas"}

	b, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("Error raised: %v", err)
		return nil, err
	}

	tmp, err := ioutil.ReadAll(bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Read error: %v", err)
		return nil, err
	}
	//fmt.Println(b, bytes.NewBuffer(b))
	err = json.Unmarshal(tmp, &spec)
	if err != nil {
		fmt.Printf("pb unmarshal error: %v", err)
		os.Exit(0)
	}
	spec.ConfidentialityLevel = pb.ConfidentialityLevel_CONFIDENTIAL
	//fmt.Println(spec2)
	return &spec, nil
}

func InvokeChaincodeSpec() (*pb.ChaincodeInvocationSpec, error) {
	var spec pb.ChaincodeSpec
	//var spec2 pb.ChaincodeSpec
	t := &params{
		1,
		map[string]string{"name": info},
		ctorMsg{"invoke", []string{"a", "b", "1"}},
		"lukas"}
	b, err := json.Marshal(t)
	if err != nil {
		panic(fmt.Errorf("Error raised: %v", err))
		return nil, err
	}

	tmp, err := ioutil.ReadAll(bytes.NewBuffer(b))
	if err != nil {
		panic(fmt.Errorf("Read error: %v", err))
		return nil, err
	}
	//fmt.Println(b, bytes.NewBuffer(b))
	err = json.Unmarshal(tmp, &spec)
	if err != nil {
		panic(fmt.Errorf("pb unmarshal error: %v", err))
		os.Exit(0)
	}
	//spec.ConfidentialityLevel = pb.ConfidentialityLevel_CONFIDENTIAL
	//fmt.Println(spec2)
	return &pb.ChaincodeInvocationSpec{&spec, ""}, nil

}



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


func CreateInvokeTx(chaincodeInvocation *pb.ChaincodeInvocationSpec, uuid string, nonce []byte, attrs ...string) (*pb.Transaction, error) {
	tx, err := pb.NewChaincodeExecute(chaincodeInvocation, uuid, pb.Transaction_CHAINCODE_INVOKE)
	if err != nil {
		fmt.Println("Failed creating new transaction")
		return nil, err
	}
	tx.Metadata, err = getMetadata(chaincodeInvocation.GetChaincodeSpec())
	if err != nil {
		fmt.Println("Failed loading Metadata")
		return nil, err
	}
	if nonce == nil {
		tx.Nonce, err = primitives.GetRandomNonce()
		if err != nil {
			fmt.Println("Failed generating Nonce")
			return nil, err
		}
	} else {
		tx.Nonce = nonce
	}
	//tx.ConfidentialityLevel = pb.ConfidentialityLevel_CONFIDENTIAL
	//tx.ConfidentialityProtocolVersion = "1.2"
	//err = encryptTx(tx)
	//if err != nil {
	//	return nil, err
	//}
	return tx, nil
}


func CreateDeployTx(chaincodeDeploymentSpec *pb.ChaincodeDeploymentSpec, uuid string, nonce []byte, attrs ...string) (*pb.Transaction, error) {
	tx, err := pb.NewChaincodeDeployTransaction(chaincodeDeploymentSpec, uuid)
	if err != nil {
		return nil, err
	}
	tx.Metadata, err = getMetadata(chaincodeDeploymentSpec.GetChaincodeSpec())
	if err != nil {
		return nil, err
	}

	if nonce == nil {
		tx.Nonce, err = primitives.GetRandomNonce()
		if err != nil {
			return nil, err
		}
	} else {
		tx.Nonce = nonce
	}
	tx.ConfidentialityLevel = pb.ConfidentialityLevel_CONFIDENTIAL
	tx.ConfidentialityProtocolVersion = "1.2"
	err = encryptTx(tx)
	if err != nil {
		return nil, err
	}

	return tx, nil

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
	Init()
	err := initViper.SetConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading viper config file"))
	}
	chaincodeInvocationSpec, err := InvokeChaincodeSpec()
	if err != nil {
		os.Exit(0)
	}
	//fmt.Println(chaincodeInvocationSpec)

	if err != nil {
		os.Exit(0)
	}

	tx, err := CreateInvokeTx(chaincodeInvocationSpec, util.GenerateUUID(), nil, chaincodeInvocationSpec.ChaincodeSpec.Attributes...)
	if err != nil {
		os.Exit(0)
	}

	tx, err = Sign(tx)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	return tx
	//fmt.Println(tx.Nonce)
	//fmt.Println(tx)
	//_ = rpc.Connect(tx)
	//fmt.Println(response)
}

func main() {
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



}



