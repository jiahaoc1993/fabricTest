package main

import (
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
	pb "github.com/hyperledger/fabric/protos"
	"tool/loadKey"
	"tool/rpc"
	"tool/initViper"
)

const (
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
func MakeAChaincodeSpec() (*pb.ChaincodeSpec, error) {
	var spec pb.ChaincodeSpec
	//var spec2 pb.ChaincodeSpec
	t := &params{
		1,
		map[string]string{"path": "github.com/hyperledger/fabric/examples/chaincode/go/chaincode_example02"},
		ctorMsg{"init", []string{"a","100000","b","10000"}},
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
		privaBytes []byte
	)

	stateKey, err = primitives.GenAESKey()
	if err != nil {
		panic(fmt.Errorf("Failed creating state key: %v\n", err))
		return err
	}

	privaBytes, err = eciesSPI.SerializePrivateKey(ccPrivateKey)
	if err != nil {
		panic(fmt.Errorf("Failed serializing chaincode key: %v\n", err))
		return err
	}

	chainPublicKey, err := loadKey.LoadKey()
	if err != nil {
		fmt.Println("error")
	}

	//fmt.Println("This is chainPublicKey: ", chainPublicKey.pub)

	cipher, err := eciesSPI.NewAsymmetricCipherFromPublicKey(chainPublicKey)

	msgToValidators, err := asn1.Marshal(chainCodeValidatorMessage1_2{privaBytes, stateKey})
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

func Sign(tx *pb.Transaction) error {
	enrollmentCert, privKey, err := loadKey.LoadEnrollment()
	if err != nil {
		fmt.Printf("Failed loading enrollment metieral")
		return err
	}

	tx.Cert = enrollmentCert.Raw

	rawTx, err := proto.Marshal(tx)
	if err != nil {
		fmt.Printf("Failed marshaling tx: %v", err)
		return err
	}

	rawSignature, err := primitives.ECDSASign(privKey, rawTx)
	if err != nil {
		fmt.Println("Failed Creating signature: %v", err)
		return err
	}

	tx.Signature = rawSignature

	return nil
}

func main() {
	Init()
	err := initViper.SetConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading viper config file"))
	}
	spec, err := MakeAChaincodeSpec()
	if err != nil {
		os.Exit(0)
	}
	fmt.Println(spec)

	chaincodeDeploymentSpec, err := getChaincodeBytes(context.Background(), spec)

	if err != nil {
		os.Exit(0)
	}

	tx, err := CreateDeployTx(chaincodeDeploymentSpec, chaincodeDeploymentSpec.ChaincodeSpec.ChaincodeID.Name, nil, spec.Attributes...)
	if err != nil {
		os.Exit(0)
	}

	err = Sign(tx)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	//fmt.Println(tx.Nonce)

	response := rpc.Connect(tx)
	fmt.Println(response)
}
