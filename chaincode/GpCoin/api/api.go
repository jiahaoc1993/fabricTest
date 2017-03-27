package main

import (
	"fmt"
	//"os"

	//"reflect"
	"time"
//	"flag"
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

func deploy() (resp string,err error) {
	appLogger.Debug("------------- Alice wants to assign the administrator role to Bob;")
	bobCert, err = bob.GetEnrollmentCertificateHandler()
	if err != nil {
		appLogger.Errorf("Failed getting Bob ECert [%s]", err)
		return
	}
	resp, err = deployInternal(bob, bobCert)
	if err != nil {
		appLogger.Errorf("Failed deploying [%s]", err)
		return
	}

	fmt.Println("chaincode name: ", resp)

	appLogger.Debug("Wait 30 seconds")
	time.Sleep(30 * time.Second)
	return
}

func topup(chaincodeName string, amount string, user string) (err error) {
	appLogger.Debug("------------- topup...")
	bobCert, err = bob.GetEnrollmentCertificateHandler()
	if err != nil {
		appLogger.Errorf("Failed getting Bob ECert [%s]", err)
		return
	}

	resp, err := topupInternal(chaincodeName, bob, bobCert, amount, user)
	if err != nil {
		appLogger.Errorf("Failed assigning ownership [%s]", err)
		return
	}
	appLogger.Debugf("Resp [%s]", resp.String())

	appLogger.Debug("Wait 30 seconds")
	//time.Sleep(30 * time.Second)


	appLogger.Debug("------------- Done!")
	return
}


func invest(chaincodeName string, amount string, user string) (err error) {
	appLogger.Debug("------------- invest...")

	bobCert, err = bob.GetEnrollmentCertificateHandler()
	if err != nil {
		appLogger.Errorf("Failed getting Bob ECert [%s]", err)
		return
	}
	resp, err := investInternal(chaincodeName, bob, bobCert, amount, user)
	if err != nil {
		appLogger.Errorf("Failed assigning ownership [%s]", err)
		return
	}
	appLogger.Debugf("Resp [%s]", resp.String())

	appLogger.Debug("Wait 30 seconds")
	//time.Sleep(30 * time.Second)


	appLogger.Debug("------------- Done!")
	return
}


func cashout(chaincodeName string, amount string, user string) (err error) {
	appLogger.Debug("------------- cashout...")

	bobCert, err = bob.GetEnrollmentCertificateHandler()
	if err != nil {
		appLogger.Errorf("Failed getting Bob ECert [%s]", err)
		return
	}
	resp, err := cashoutInternal(chaincodeName, bob, bobCert, amount, user)
	if err != nil {
		appLogger.Errorf("Failed assigning ownership [%s]", err)
		return
	}
	appLogger.Debugf("Resp [%s]", resp.String())

	appLogger.Debug("Wait 30 seconds")
	//time.Sleep(30 * time.Second)


	appLogger.Debug("------------- Done!")
	return
}

func transfer(chaincodeName string, amount string, from string, to string) (err error) {
	appLogger.Debug("------------- transfer...")


	bobCert, err = bob.GetEnrollmentCertificateHandler()
	if err != nil {
		appLogger.Errorf("Failed getting Bob ECert [%s]", err)
		return
	}
	resp, err := transferInternal(chaincodeName, bob, bobCert, amount, from, to)
	if err != nil {
		appLogger.Errorf("Failed assigning ownership [%s]", err)
		return
	}
	appLogger.Debugf("Resp [%s]", resp.String())

	appLogger.Debug("Wait 30 seconds")
//	time.Sleep(30 * time.Second)


	appLogger.Debug("------------- Done!")
	return
}
/*
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

*/

func check(n int, args []string) error{
	if len(args) != n {
		err := fmt.Errorf("wrong args\n")
		return err
	}

	return nil
}

/*
func main() {
	var chaincodeName   string
	var method	    string
	flag.StringVar(&chaincodeName, "n", " ", "Chaincode Name")
	flag.StringVar(&method, "m", " ", "method it call")
	flag.Parse()
	if method != "deploy" && len(chaincodeName) == 1 {
		panic(fmt.Errorf("no caincodeName\n"))

	}
	// Initialize a non-validating peer whose role is to submit
	// transactions to the fabric network.
	// A 'core.yaml' file is assumed to be available in the working directory.
	if err := initNVP(); err != nil {
		appLogger.Debugf("Failed initiliazing NVP [%s]", err)
		os.Exit(-1)
	}
	// Enable fabric 'confidentiality'
	confidentiality(false)
	args := flag.Args()
	switch method {
		case "deploy" :
				chainName := deploy()
				fmt.Println(chainName)
		case "topup"  :
				check(2, args)
				topup(chaincodeName, args[0], args[1])
		case "invest" :
				check(2, args)
				invest(chaincodeName, args[0], args[1])
		case "cashout":
				check(2, args)
				cashout(chaincodeName, args[0], args[1])
		case "transfer":
				check(3, args)
				transfer(chaincodeName, args[0], args[1], args[2])
		case "query":
				check(1, args)
				CheckUser(chaincodeName, args[0])
		default :
			fmt.Println("you must input a method")

	}
*/

	// Exercise the 'asset_management' chaincode
	//if err := testGPCoinChaincode(); err != nil {
	//	appLogger.Debugf("Failed testing asset management chaincode [%s]", err)
	//	os.Exit(-2)
	//}

