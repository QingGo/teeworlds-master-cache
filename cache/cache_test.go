package cache

import (
	"fmt"
	"testing"
)

func TestGetServerInfoFromRestAPI(t *testing.T) {
	getServerInfoFromRestAPI()
	fmt.Println(ServerAddrList)
}

func TestGetServerInfoFromMaster(t *testing.T) {
	getServerInfoFromMasterList()
	fmt.Println(ServerAddrList)
}
