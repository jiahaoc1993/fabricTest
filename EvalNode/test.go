package main

import (
	"strings"
	"fmt"
	"github.com/spf13/viper"
	//"tool/loadKey"
	//pb "github.com/hyperledger/fabric/protos"
	//"tool/rpc"
)

func main() {
	viper.SetEnvPrefix("core")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AddConfigPath("/opt/gopath/src/github.com/hyperledger/fabric/peer/")
	viper.SetConfigName("core")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error when reading config file"))
	}

	fmt.Println(viper.GetInt("peer.gomaxprocs"))
}
