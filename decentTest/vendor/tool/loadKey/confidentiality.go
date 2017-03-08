package loadKey

import (
	"crypto/ecdsa"
	//	"crypto/rand"
	"crypto/x509"
	"fmt"
	"github.com/hyperledger/fabric/core/crypto/primitives"
	"github.com/hyperledger/fabric/core/crypto/primitives/ecies"
	"io/ioutil"
	//	"os"
	//	"io"
)

func loadEnrollmentKey() (*ecdsa.PrivateKey, error) {
	path := "/var/hyperledger/production/crypto/client/diego/ks/raw/enrollment.key"
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
		//panic(fmt.Errorf("Failed loading private key: %v", err))
	}

	privateKey, err := primitives.PEMtoPrivateKey(raw, nil)
	if err != nil {
		return nil, err
		//panic(fmt.Errorf("Failed parsing private key: %v", err))
	}
	return privateKey.(*ecdsa.PrivateKey), nil

}

func loadFakeEnrollmentKey() (*ecdsa.PrivateKey, error) {
	path := "/var/hyperledger/production/crypto/client/diego/ks/raw/enrollment2.key"
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
		//panic(fmt.Errorf("Failed loading private key: %v", err))
	}

	privateKey, err := primitives.PEMtoPrivateKey(raw, nil)
	if err != nil {
		return nil, err
		//panic(fmt.Errorf("Failed parsing private key: %v", err))
	}
	return privateKey.(*ecdsa.PrivateKey), nil

}

func loadEnrollmentCertificate() (*x509.Certificate, []byte, error) {
	path := "/var/hyperledger/production/crypto/client/diego/ks/raw/enrollment.cert"
	pem, err := ioutil.ReadFile(path)
	if err != nil {
		//panic(fmt.Errorf("Failed loading certificate: %v", err))
		return nil, nil, err
	}
	cert, der, err := primitives.PEMtoCertificateAndDER(pem)
	if err != nil {
		return nil, nil, err
		//fmt.Errorf("Failed parsing certificate: %v", err))
	}
	return cert, der, nil
}

func LoadEnrollment() (*x509.Certificate, *ecdsa.PrivateKey, error) {
	cert, _, err := loadEnrollmentCertificate()
	if err != nil {
		fmt.Printf("loadEnrollmentCertificate failed: %v\n", err)
		return nil, nil, err
	}
	pk := cert.PublicKey.(*ecdsa.PublicKey) // public key in enrollment certificate

	privKey, err := loadEnrollmentKey() // private key
	if err != nil {
		fmt.Printf("loadEnrollmentKey failed: %v\n", err)
		return nil, nil, err
	}

	err = primitives.VerifySignCapability(privKey, pk)
	if err != nil {
		fmt.Println("Failed Checking enrollment certificate against enrollment key: %v", err.Error())
		return nil, nil, err
	}
	return cert, privKey, nil
}

func LoadFakeEnrollment() (*x509.Certificate, *ecdsa.PrivateKey, error) {
	cert, _, err := loadEnrollmentCertificate()
	if err != nil {
		fmt.Printf("loadEnrollmentCertificate failed: %v\n", err)
		return nil, nil, err
	}

	privKey, err := loadFakeEnrollmentKey() // private key
	if err != nil {
		fmt.Printf("loadEnrollmentKey failed: %v\n", err)
		return nil, nil, err
	}

	return cert, privKey, nil
}

func LoadKey() (primitives.PublicKey, error) {
	eciesSPI := ecies.NewSPI()

	pathChainKey := "/var/hyperledger/production/crypto/client/diego/ks/raw/chain.key"
	raw, err := ioutil.ReadFile(pathChainKey)
	if err != nil {
		panic(fmt.Errorf("Failed loading private key: %v\n", err))
		return nil, err
	}
	//fmt.Println(raw)
	publicKey, err := primitives.PEMtoPublicKey(raw, nil)
	if err != nil {
		panic(fmt.Errorf("Failed parsing private key: %v\n", err))
		return nil, err

	}
	//fmt.Println(publicKey)
	t, ok := publicKey.(*ecdsa.PublicKey)
	if ok {
		//fmt.Println("1")
		chainPublicKey, err := eciesSPI.NewPublicKey(nil, t)
		if err != nil {
			fmt.Println("Wrong")
			return nil, err
		}
		return chainPublicKey, nil
		/*
			cipher, err := eciesSPI.NewAsymmetricCipherFromPublicKey(chainPublicKey)
			if err != nil {
				panic(fmt.Errorf("Failed creating new encryption shcheme: %v", err))
				return nil, err
			}
			return cipher, nil
		*/
	}
	return nil, err
}
