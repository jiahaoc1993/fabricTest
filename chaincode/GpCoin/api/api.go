package main

import (
	//"fmt"
	"os"
	//"reflect"
	"time"

	"github.com/hyperledger/fabric/core/crypto"
	pb "github.com/hyperledger/fabric/protos"
	"github.com/op/go-logging"
	"google.golang.org/grpc"
)

var (
	// Logging
	appLogger = logging.MustGetLogger("app")

	// NVP related objects
	peerClientConn *grpc.ClientConn
	serverClient   pb.PeerClient

	// Alice is the deployer
	alice crypto.Client
	aliceCert crypto.CertificateHandler
	// Bob is the administrator
	bob     crypto.Client
	bobCert crypto.CertificateHandler

	// Charlie and Dave are owners
	charlie     crypto.Client
	charlieCert crypto.CertificateHandler

	dave     crypto.Client
	daveCert crypto.CertificateHandler
)

func deploy() (err error) {
	appLogger.Debug("------------- Alice wants to assign the administrator role to Bob;")
	// Deploy:
	// 1. Alice is the deployer of the chaincode;
	// 2. Alice wants to assign the administrator role to Bob;
	// 3. Alice obtains, via an out-of-band channel, a TCert of Bob, let us call this certificate *BobCert*;

	bobCert, err = bob.GetTCertificateHandlerNext()
	if err != nil {
		appLogger.Errorf("Failed getting Bob TCert [%s]", err)
		return
	}

	// 4. Alice constructs a deploy transaction, as described in *application-ACL.md*,  setting the transaction
	// metadata to *DER(CharlieCert)*.
	// 5. Alice submits th	e transaction to the fabric network.
	resp, err := deployInternal(alice, bobCert)
	if err != nil {
		appLogger.Errorf("Failed deploying [%s]", err)
		return
	}
	appLogger.Debugf("Resp [%s]", resp.String())
	appLogger.Debugf("Chaincode NAME: [%s]-[%s]", chaincodeName, string(resp.Msg))

	appLogger.Debug("Wait 30 seconds")
	time.Sleep(30 * time.Second)

	appLogger.Debug("------------- Done!")
	return
}

func topup() (err error) {
	appLogger.Debug("------------- topup...")
	//charlie topUp 10000 USD
	charlieCert, err = charlie.GetTCertificateHandlerNext()
	if err != nil {
		appLogger.Errorf("Failed getting Charlie TCert [%s]", err)
		return
	}

	resp, err := topupInternal(bob, bobCert, "10000", charlieCert)
	if err != nil {
		appLogger.Errorf("Failed assigning ownership [%s]", err)
		return
	}
	appLogger.Debugf("Resp [%s]", resp.String())

	appLogger.Debug("Wait 30 seconds")
	time.Sleep(30 * time.Second)


	appLogger.Debug("------------- Done!")
	return
}


func invest() (err error) {
	appLogger.Debug("------------- invest...")
	//charlie topUp 100 USD
	charlieCert, err = charlie.GetTCertificateHandlerNext()
	if err != nil {
		appLogger.Errorf("Failed getting Charlie TCert [%s]", err)
		return
	}

	resp, err := investInternal(bob, bobCert, "1000", charlieCert)
	if err != nil {
		appLogger.Errorf("Failed assigning ownership [%s]", err)
		return
	}
	appLogger.Debugf("Resp [%s]", resp.String())

	appLogger.Debug("Wait 30 seconds")
	time.Sleep(30 * time.Second)


	appLogger.Debug("------------- Done!")
	return
}


func cashout() (err error) {
	appLogger.Debug("------------- cashout...")
	//charlie topUp 10000 USD
	charlieCert, err = charlie.GetTCertificateHandlerNext()
	if err != nil {
		appLogger.Errorf("Failed getting Charlie TCert [%s]", err)
		return
	}

	resp, err := topupInternal(bob, bobCert, "0.8", charlieCert)
	if err != nil {
		appLogger.Errorf("Failed assigning ownership [%s]", err)
		return
	}
	appLogger.Debugf("Resp [%s]", resp.String())

	appLogger.Debug("Wait 30 seconds")
	time.Sleep(30 * time.Second)


	appLogger.Debug("------------- Done!")
	return
}

func transfer() (err error) {
	appLogger.Debug("------------- topup...")
	//charlie topUp 10000 USD
	charlieCert, err = charlie.GetTCertificateHandlerNext()
	if err != nil {
		appLogger.Errorf("Failed getting Charlie TCert [%s]", err)
		return
	}

	aliceCert, err = alice.GetTCertificateHandlerNext()
	if err != nil {
		appLogger.Errorf("Failed getting Charlie TCert [%s]", err)
		return
	}

	resp, err := transferInternal(bob, bobCert, "1.5", charlieCert, aliceCert)
	if err != nil {
		appLogger.Errorf("Failed assigning ownership [%s]", err)
		return
	}
	appLogger.Debugf("Resp [%s]", resp.String())

	appLogger.Debug("Wait 30 seconds")
	time.Sleep(30 * time.Second)


	appLogger.Debug("------------- Done!")
	return
}

func testGPCoinChaincode() (err error) {
	// Deploy
	err = deploy()
	if err != nil {
		appLogger.Errorf("Failed deploying [%s]", err)
		return
	}

	// topup
	err = topup()
	if err != nil {
		appLogger.Errorf("Failed assigning ownership [%s]", err)
		return
	}

	// invest
	err = invest()
	if err != nil {
		appLogger.Errorf("Failed transfering ownership [%s]", err)
		return
	}

	err = cashout()
	if err != nil {
		appLogger.Errorf("Failed transfering ownership [%s]", err)
		return
	}

	err = transfer()
	if err != nil {
		appLogger.Errorf("Failed transfering ownership [%s]", err)
		return
	}

	appLogger.Debug("Dave is the owner!")

	return
}

func main() {
	// Initialize a non-validating peer whose role is to submit
	// transactions to the fabric network.
	// A 'core.yaml' file is assumed to be available in the working directory.
	if err := initNVP(); err != nil {
		appLogger.Debugf("Failed initiliazing NVP [%s]", err)
		os.Exit(-1)
	}

	// Enable fabric 'confidentiality'
	confidentiality(true)

	// Exercise the 'asset_management' chaincode
	if err := testGPCoinChaincode(); err != nil {
		appLogger.Debugf("Failed testing asset management chaincode [%s]", err)
		os.Exit(-2)
	}
}

