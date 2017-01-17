package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/hyperledger/fabric/core/crypto/primitives/ecies"
	"io/ioutil"
	//	"os"
	"io"
)

type publicKeyImpl struct {
	pub    *ecdsa.PublicKey
	rand   io.Reader
	params *ecies.Params
}

func (pk *publicKeyImpl) GetRand() io.Reader {
	return pk.rand
}

func (pk *publicKeyImpl) IsPublic() bool {
	return true
}

func newPublicKeyFromECDSA(r io.Reader, pk *ecdsa.PublicKey) (primitives.PublicKey, error) {
	if r == nil {
		r = rand.Reader
	}

	if pk == nil {
		return nil, fmt.Errorf("Null ECDSA public key")
	}

	return &publicKeyImpl{pk, r, nil}, nil
}

func main() {
	eciesSPI := ecies.NewSPI()
	path := "/var/hyperledger/production/crypto/client/diego/ks/raw/query.key"
	pem, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("Read Error: %v\n", err))
	}
	fmt.Println(pem)

	pathChainKey := "/var/hyperledger/production/crypto/client/diego/ks/raw/chain.key"
	raw, err := ioutil.ReadFile(pathChainKey)
	if err != nil {
		panic(fmt.Errorf("Failed loading private key: %v\n", err))
	}
	fmt.Println(raw)
	publicKey, err := primitives.PEMtoPublicKey(raw, nil)
	if err != nil {
		panic(fmt.Errorf("Failed parsing private key: %v\n", err))

	}
	fmt.Println(publicKey)
	t, ok := publicKey.(*ecdsa.PublicKey)
	if ok {
		//fmt.Println("1")
		chainPublicKey, err := eciesSPI.NewPublicKey(nil, t)
		if err != nil {
			fmt.Println("Wrong")
		}

		cipher, err := eciesSPI.NewAsymmetricCipherFromPublicKey(chainPublicKey)
		if err != nil {
			panic(fmt.Errorf("Failed creating new encryption shcheme: %v", err))
		}
		fmt.Println(cipher)

	}

}
