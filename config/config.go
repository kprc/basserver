package config

import (
	"encoding/json"
	"github.com/kprc/nbsnetwork/tools"
	"log"
	"sync"
	"path"
	"os"
)

const(
	BASD_HomeDir = ".basd"
	BASD_CFG_FileName = "basd.json"
)

type BASDConfig struct {
	UpdPort    int    `json:"updport"`
	TcpPort    int    `json:"tcpport"`
	RopstenNAP string `json:"ropstennap"`
	TokenAddr  string `json:"tokenaddr"`
	MgrAddr    string `json:"mgraddr"`
	CmdListenPort string `json:"cmdlistenport"`
}

var (
	bascfgInst     *BASDConfig
	bascfgInstLock sync.Mutex
)

func (bc *BASDConfig) InitCfg() *BASDConfig{
	bc.UpdPort = 53
	bc.TcpPort = 53
	bc.CmdListenPort = "127.0.0.1:59527"

	return bc
}

func (bc *BASDConfig)Load()  *BASDConfig{
	if !tools.FileExists(GetBASDCFGFile()){
		return nil
	}

	jbytes,err:=tools.OpenAndReadAll(GetBASDCFGFile())
	if err!=nil{
		log.Println("load file failed",err)
		return nil
	}

	//bc1:=&BASDConfig{}

	err = json.Unmarshal(jbytes,bc)
	if err!=nil{
		log.Println("load configuration unmarshal failed",err)
		return nil
	}

	return bc

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

func PreLoad() *BASDConfig {
	bc:=&BASDConfig{}

	return bc.Load()
}


func LoadFromCfgFile(file string) *BASDConfig {
	bc := &BASDConfig{}

	bc.InitCfg()

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

func LoadFromCmd(initfromcmd func(cmdbc *BASDConfig) *BASDConfig ) *BASDConfig {
	bascfgInstLock.Lock()
	defer bascfgInstLock.Unlock()

	lbc:= newBasDCfg().Load()

	if lbc !=nil{
		bascfgInst = lbc
	}else{
		lbc=newBasDCfg()
	}

	bascfgInst = initfromcmd(lbc)

	return bascfgInst
}

func GetBASDHomeDir() string {
	curHome, err:= tools.Home()
	if err!=nil{
		log.Fatal(err)
	}

	return path.Join(curHome,BASD_HomeDir)
}

func GetBASDCFGFile() string {
	return path.Join(GetBASDHomeDir(),BASD_CFG_FileName)
}

func (bc *BASDConfig)Save()  {
	jbytes,err:=json.MarshalIndent(*bc," ","\t")

	if err!=nil{
		log.Println("Save BASD Configuration json marshal failed",err)
	}

	if !tools.FileExists(GetBASDHomeDir()){
		os.MkdirAll(GetBASDHomeDir(),0755)
	}

	err = tools.Save2File(jbytes,GetBASDCFGFile())
	if err!=nil{
		log.Println("Save BASD Configuration to file failed",err)
	}

}

func IsInitialized()  bool{
	if tools.FileExists(GetBASDCFGFile()){
		return true
	}

	return false
}