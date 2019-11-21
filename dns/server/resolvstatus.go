package server

import (
	"sync"
	"github.com/kprc/basserver/config"
	"sort"
	"fmt"
)

const(
	ResolvReuseTime int64 = 86400000   //ms
	ResolvNormal int32 = 1
	ResolvAbnormal int32 = 2
)


type ResolvStatus struct {
	Status int32
	FailCnt int
	IPStr  string
	Idx int
}

func (rs *ResolvStatus)String() string {
	return fmt.Sprint("ip: ",rs.IPStr,"  Status: ",rs.Status,"  failcnt: ",rs.FailCnt,"  idx: ",rs.Idx)
}


var (
	gResolvStatusArr []*ResolvStatus
	gResolvStatusArrLock sync.Mutex

)

func newResolvStatus(dns string,idx int) *ResolvStatus {
	return &ResolvStatus{Status:ResolvNormal,IPStr:dns,Idx:idx}
}

func getResolvStatusArr() []*ResolvStatus {
	if gResolvStatusArr == nil{
		gResolvStatusArrLock.Lock()
		defer gResolvStatusArrLock.Unlock()

		if gResolvStatusArr == nil{
			cfg:=config.GetBasDCfg()
			gResolvStatusArr = make([]*ResolvStatus,len(cfg.ResolvDns))
			for idx,dns:=range cfg.ResolvDns{
				gResolvStatusArr[idx] = newResolvStatus(dns,idx)
			}
		}
	}
	return gResolvStatusArr
}

func GetDns() string  {
	ndns := getResolvStatusArr()

	gResolvStatusArrLock.Lock()
	defer  gResolvStatusArrLock.Unlock()

	//fmt.Println("GetDns CurIdx: ",curIdx, " IP: ",ndns[curIdx].IPStr)
	//
	//fmt.Println("=========")
	//for _,n:=range ndns{
	//	fmt.Println(n.String())
	//}
	//fmt.Println("=========")

	return ndns[0].IPStr

}

func FailDns(ips string)  {
	ndns:=getResolvStatusArr()

	gResolvStatusArrLock.Lock()
	defer  gResolvStatusArrLock.Unlock()

	for i:=0;i<len(ndns);i++{
		if ndns[i].IPStr == ips{
			ndns[i].FailCnt ++
		}
	}

	sort.Slice(ndns, func(i, j int) bool {
		if ndns[i].FailCnt < ndns[j].FailCnt{
			return true
		}
		return false
	})

	//for _,n:=range ndns{
	//	fmt.Println(n.String())
	//}
	//fmt.Println("curIdx:",curIdx)

	//curIdx = ndns[0].Idx

}

