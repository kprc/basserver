package config

import "sync"

type BASDConfig struct {
	UpdPort int
	TcpPort int

}


var(
	bascfgInst *BASDConfig
	bascfgInstLock sync.Mutex
)

func (bc *BASDConfig)InitCfg()  {
	bc.UpdPort = 53
	bc.TcpPort = 53
}

func newBasDCfg() *BASDConfig  {

	bc:=&BASDConfig{}

	bc.InitCfg()

	return bc
}

func GetBasDCfg() *BASDConfig {
	if bascfgInst == nil{
		bascfgInstLock.Lock()
		defer bascfgInstLock.Unlock()
		if bascfgInst == nil{
			bascfgInst = newBasDCfg()
		}
	}

	return bascfgInst
}


