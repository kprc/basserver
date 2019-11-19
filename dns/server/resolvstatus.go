package server

import (
	"sync"
	"github.com/kprc/basserver/config"
	"github.com/kprc/nbsnetwork/tools"
)

const(
	ResolvReuseTime int64 = 86400000   //ms
	ResolvNormal int32 = 1
	ResolvAbnormal int32 = 2
)


type ResolvStatus struct {
	Status int32
	LastFailTime int64
	IPStr  string
}


var (
	gResolvStatusArr []*ResolvStatus
	gResolvStatusArrLock sync.Mutex
)

func newResolvStatus(dns string) *ResolvStatus {
	return &ResolvStatus{Status:ResolvNormal,IPStr:dns}
}

func getResolvStatusArr() []*ResolvStatus {
	if gResolvStatusArr == nil{
		gResolvStatusArrLock.Lock()
		defer gResolvStatusArrLock.Unlock()

		if gResolvStatusArr == nil{
			cfg:=config.GetBasDCfg()
			gResolvStatusArr = make([]*ResolvStatus,len(cfg.ResolvDns))
			for idx,dns:=range cfg.ResolvDns{
				gResolvStatusArr[idx] = newResolvStatus(dns)
			}
		}
	}
	return gResolvStatusArr
}

func GetDns() string  {
	ndns := getResolvStatusArr()

	gResolvStatusArrLock.Lock()
	defer  gResolvStatusArrLock.Unlock()

	for i:=0;i<len(ndns);i++{
		if ndns[i].LastFailTime == 0{
			return ndns[i].IPStr
		}
	}

	now := tools.GetNowMsTime()

	for i:=0; i<len(ndns);i++{
		if now - ndns[i].LastFailTime > ResolvReuseTime{
			ndns[i].LastFailTime = 0
			ndns[i].Status = ResolvNormal

			return  ndns[i].IPStr
		}
	}

	return ""

}

func FailDns(ips string)  {
	ndns := getResolvStatusArr()

	gResolvStatusArrLock.Lock()
	defer  gResolvStatusArrLock.Unlock()

	for i:=0;i<len(ndns);i++{
		if ndns[i].IPStr == ips{
			ndns[i].Status = ResolvAbnormal
			ndns[i].LastFailTime = tools.GetNowMsTime()
		}
	}
}