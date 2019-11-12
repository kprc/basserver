package config

import (
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"sync"
)

type BASDConfig struct {
	UpdPort    int    `json:"updport"`
	TcpPort    int    `json:"tcpport"`
	RopstenNAP string `json:"ropstennap"`
	TokenAddr  string `json:"tokenaddr"`
	MgrAddr    string `json:"mgraddr"`
}

var (
	bascfgInst     *BASDConfig
	bascfgInstLock sync.Mutex
)

func (bc *BASDConfig) InitCfg() {
	bc.UpdPort = 53
	bc.TcpPort = 53
}

func newBasDCfg() *BASDConfig {

	bc := &BASDConfig{}

	bc.InitCfg()

	return bc
}

func GetBasDCfg() *BASDConfig {
	if bascfgInst == nil {
		bascfgInstLock.Lock()
		defer bascfgInstLock.Unlock()
		if bascfgInst == nil {
			bascfgInst = newBasDCfg()
		}
	}

	return bascfgInst
}

func LoadFromCfgFile(file string) *BASDConfig {
	bc := &BASDConfig{}

	bcontent, err := tools.OpenAndReadAll(file)
	if err != nil {
		log.Fatal("Load Config file failed")
		return nil
	}

	err = json.Unmarshal(bcontent, bc)
	if err != nil {
		log.Fatal("Load Config From json failed")
		return nil
	}

	bascfgInstLock.Lock()
	defer bascfgInstLock.Unlock()
	bascfgInst = bc

	return bc

}

func LoadFromCmd(bc *BASDConfig) *BASDConfig {
	bascfgInstLock.Lock()
	defer bascfgInstLock.Unlock()

	bascfgInst = bc

	return bascfgInst
}
