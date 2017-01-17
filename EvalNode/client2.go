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
	"golang.org/x/net/context"
	//"github.com/golang/protobuf/jsonpb"
	//"github.com/hyperledger/fabric/core/crypto"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/hyperledger/fabric/core/crypto/primitives/ecies"
	pb "github.com/hyperledger/fabric/protos"
	"tool/loadKey"
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
		map[string]string{"path": "github.com/hyperledger/fabric/examples/chaincode/go/HelloWorld"},
		ctorMsg{"init", []string{"Hello, Pig"}},
		"diego"}

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
	//handle confidentiality
	//fmt.Println(chaincodeDeploymentSpec.ChaincodeSpec.ConfidentialityLevel)
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

	_, err = asn1.Marshal(chainCodeValidatorMessage1_2{privaBytes, stateKey})
	if err != nil {
		panic(fmt.Errorf("Failed to preparing message to the validators: %v", err))
	}

	encMsgToValidators, err := cipher.Process([]byte{123})
	if err != nil {
		panic(fmt.Errorf("Failed to encrypting message to the validators: %v", err))
	}
	fmt.Println(encMsgToValidators)

	return nil

}

func main() {
	Init()
	//configuration
	//for viper testing
	/*
		viper.SetEnvPrefix("core")
		viper.AutomaticEnv()
		replacer := strings.NewReplacer(".", "_")
		viper.SetEnvKeyReplacer(replacer)
		viper.SetConfigName("core")
		viper.AddConfigPath("/opt/gopath/src/github.com/hyperledger/fabric/peer/")
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("error raise: %v", err))
		}
		//viper.Set("peer.fileSystemPath", filepath.Join("/", "var", "hyperledger", "production"))
		err = core.CacheConfiguration()
		if err != nil {
			panic(fmt.Errorf("error raise: %v", err))
		}
		//viper.AddConfigPath("/hyperledger/peer/")
		//viper.AddConfigPath("/opt/gopath/src/github.com/hyperledger/fabric/peer/")
		//fmt.Println(viper.GetString("peer.fileSystemPath"))
		//fmt.Println(viper.GetString("peer.gomaxprocs"))
		fmt.Println(viper.GetBool("security.enabled"))
		fmt.Println(peer.SecurityEnabled())
		//	fmt.Println(string(viper.GetString("peer.validator.consensus.plugin")))

		//fmt.Println(viper.GetString("chaincode.mode") == chaincode.DevModeUserRunsChaincode)
		//fmt.Println(viper.GetBool("security.privacy"))
		//fmt.Println(viper.GetBool("security.enabled"))

		//define the devop server
		//var serverDevops pb.DevopsServer
		//serverDevops = //use underlying *core.Devops
		//var spec pb.ChaincodeSpec
		//	t, err := MakeATransaction()
		//transId := new(string)
	*/
	spec, err := MakeAChaincodeSpec()
	if err != nil {
		os.Exit(0)
	}
	fmt.Println(spec)

	chaincodeDeploymentSpec, err := getChaincodeBytes(context.Background(), spec)

	if err != nil {
		os.Exit(0)
	}

	_, err = CreateDeployTx(chaincodeDeploymentSpec, chaincodeDeploymentSpec.ChaincodeSpec.ChaincodeID.Name, []byte{}, spec.Attributes...)
	if err != nil {
		os.Exit(0)
	}

	/*
		transId, err = Deploy(context.Background(), spec)

		if err != nil {
			fmt.Printf("Error raised: %v", err)
			os.Exit(0)
		}

		fmt.Println(transId)
	*/
	/*
		chaincodeDeploymentSpec, err := getChaincodeBytes(context.Background(), spec)
		if err != nil {
			fmt.Printf("error raised: %v", err)
		}
	*/
	//fmt.Println(chaincodeDeploymentSpec.ChaincodeSpec.ChaincodeID.Name)
}
