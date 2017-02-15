package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	pb "github.com/hyperledger/fabric/protos"
	"github.com/hyperledger/fabric/core/crypto/primitives"
)

const (
	//info string = "1600b975353e233708899a3b0ff8da55418f0738ef47f4a22d84b90da481d31261432209f1e7c4767dd2c400d685f4c96d41493e4576a52e41aaa36b142eaf81"
	addr string = "http://10.0.2.15:7050"
	//info string = "0950f1683d36ebd721845444b9fd1b6e732bdc12650e2b0d9ea77f4a7b83856dd830e938ed6473f7398b2bb47f6b7a80dfbc12222409b9d8365d0d8f73064535"
//	info string = "fcf53ff7f2a6204a4ec3d50943d7b8101853d5446a004659418b3c0621b1409fb24b825901ecc53b5927b535faca423098f88c1323eccf1bad658a9ad45a32ff"
)


type params struct {
	Type		int		  `json:"type,omitempty"`
	ChaincodeID	map[string]string `json:"chaincodeID,omitempty"`
	CtorMsg		ctorMsg		  `json:"CtorMsg"`
	SecureContext	string		  `json:"secureContext,omitempty"`
}

type ctorMsg struct {
	Function string		`json:"function,omitempty"`
	Args	 []string	`json:"args,omitempty"`
}

func InvokeChaincodeSpec(info string) (*pb.ChaincodeInvocationSpec, error) {
	var spec pb.ChaincodeSpec
	t := &params{
		1,
		map[string]string{"name": info},
		ctorMsg{"invoke", []string{"a", "b", "1"}},
		"lukas"}
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("Error raised: %v\n", err)
		return nil, err
	}

	tmp, err := ioutil.ReadAll(bytes.NewBuffer(b)) //read the transmitted json bytes
	if err != nil {
		fmt.Printf("Error raised: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(tmp, &spec)
	if err != nil {
		fmt.Printf("Failed unmarshaling json to spec: %v\n", err)
	}
	return &pb.ChaincodeInvocationSpec{&spec, ""}, nil
}


func QueryChaincodeSpec(info string) (*pb.ChaincodeInvocationSpec, error) {
	var spec pb.ChaincodeSpec
	t := &params{
		1,
		map[string]string{"name": info},
		ctorMsg{"query", []string{"b"}},
		"lukas"}
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("Error raised: %v\n", err)
		return nil, err
	}

	tmp, err := ioutil.ReadAll(bytes.NewBuffer(b)) //read the transmitted json bytes
	if err != nil {
		fmt.Printf("Error raised: %v\n", err)
		return nil, err
	}
	err = json.Unmarshal(tmp, &spec)
	if err != nil {
		fmt.Printf("Failed unmarshaling json to spec: %v\n", err)
	}
	return &pb.ChaincodeInvocationSpec{&spec, ""}, nil
}

//sperate create the tx!

func getMetadata(chaincodespec *pb.ChaincodeSpec) ([]byte, error) {
	return chaincodespec.Metadata, nil
}



func CreateInvokeTx(chaincodeInvocation *pb.ChaincodeInvocationSpec, uuid string, nonce []byte, attrs ...string) (*pb.Transaction, error) {
	tx, err := pb.NewChaincodeExecute(chaincodeInvocation, uuid, pb.Transaction_CHAINCODE_INVOKE)
	if err != nil {
		fmt.Printf("Error raised: %v\n", err)
		return nil, err

	}
	tx.Metadata, err = getMetadata(chaincodeInvocation.GetChaincodeSpec())
	if err != nil {
		fmt.Printf("Failed loading Metadata: %v\n", err)
		return nil, err
	}

	if nonce == nil {
		tx.Nonce, err = primitives.GetRandomNonce()
		if err != nil {
			fmt.Printf("Failed generating nonce: %v\n", err)
			return nil, err
		}
	} else {
		tx.Nonce = nonce
	}

	return tx, nil
}


func CreateQueryTx(chaincodeInvocation *pb.ChaincodeInvocationSpec, uuid string, nonce []byte, attrs ...string) (*pb.Transaction, error) {
	tx, err := pb.NewChaincodeExecute(chaincodeInvocation, uuid, pb.Transaction_CHAINCODE_QUERY)
	if err != nil {
		fmt.Printf("Error raised: %v\n", err)
		return nil, err

	}
	tx.Metadata, err = getMetadata(chaincodeInvocation.GetChaincodeSpec())
	if err != nil {
		fmt.Printf("Failed loading Metadata: %v\n", err)
		return nil, err
	}

	if nonce == nil {
		tx.Nonce, err = primitives.GetRandomNonce()
		if err != nil {
			fmt.Printf("Failed generating nonce: %v\n", err)
			return nil, err
		}
	} else {
		tx.Nonce = nonce
	}

	return tx, nil
}
