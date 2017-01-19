package initViper

import (
	"strings"
	"fmt"
	"github.com/spf13/viper"
	//"tool/loadKey"
	//pb "github.com/hyperledger/fabric/protos"
	//"tool/rpc"
)

func SetConfig() error{
	viper.SetEnvPrefix("core")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AddConfigPath("/opt/gopath/src/github.com/hyperledger/fabric/peer/")
	viper.SetConfigName("core")
	err := viper.ReadInConfig()
	if err != nil {
		//panic(fmt.Errorf("Fatal error when reading config file"))
		return  err
	}

	//fmt.Println(viper.GetInt("peer.gomaxprocs"))
	return nil
}
