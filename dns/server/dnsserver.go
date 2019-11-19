package server

import (
	"github.com/Ungigdu/BAS_contract_go/BAS_Ethereum"
	"github.com/kprc/basserver/config"
	"github.com/miekg/dns"
	"log"
	"net"
	"strconv"
	"encoding/binary"
	"github.com/btcsuite/btcutil/base58"
)

const (
	TypeBCAddr = 65
)

var (
	dnshandle dns.HandlerFunc
)

var(
	udpServer *dns.Server
	tcpServer *dns.Server
)


type DR struct {
	BAS_Ethereum.DomainRecord
}

func (dr *DR)IntIPv4() uint32 {
	 return  binary.BigEndian.Uint32(dr.IPv4[:])
}

func DnsHandle(writer dns.ResponseWriter, msg *dns.Msg) {
	//todo...
	m := msg.Copy()

	m.Compress = true
	m.Response = true

	q:=m.Question[0]

	qn := q.Name

	//log.Println("Get query string",qn)

	if qn[len(qn)-1] == '.' {
		qn = qn[:len(qn)-1]

	}

	var bdr BAS_Ethereum.DomainRecord
	var err error

	if q.Qtype == dns.TypeA{
		bdr, err = BAS_Ethereum.QueryByString(qn)
	}

	log.Println("QType: ",q.Qtype,"query string",qn)

	if q.Qtype == TypeBCAddr{
		//log.Println(qn)
		var b []byte
		b = base58.Decode(qn)
		var barr [32]byte

		for i:=0;i<len(b);i++{
			barr[i] = b[i]
		}

		bdr,err = BAS_Ethereum.QueryByBCAddress(barr)

	}

	dr := &DR{bdr}
	if err != nil || dr.IntIPv4() == 0 {
		log.Println("Can't Get Domain Name info",q.Name)
		m.Rcode = dns.RcodeBadKey

		writer.WriteMsg(m)

		return
	}

	A := &dns.A{}

	A.Hdr.Name = q.Name
	A.Hdr.Rrtype = dns.TypeA
	A.Hdr.Class = dns.ClassINET
	A.Hdr.Ttl = 10
	A.Hdr.Rdlength = 4

	A.A = net.IPv4(dr.IPv4[0], dr.IPv4[1], dr.IPv4[2], dr.IPv4[3])

	log.Println("Request Name: ", qn, A.A.String())

	var rr []dns.RR

	rr = append(rr, A)

	m.Answer = rr

	writer.WriteMsg(m)

}

func DNSServerDaemon() {
	cfg := config.GetBasDCfg()

	//fmt.Println(*cfg)

	uport := cfg.UpdPort
	uaddr := ":" + strconv.Itoa(uport)

	dnshandle = DnsHandle

	log.Println("DNS Server Start at udp", uaddr)

	udpServer = &dns.Server{}
	udpServer.Addr = uaddr
	udpServer.Handler = dnshandle
	udpServer.Net = "udp4"

	go udpServer.ListenAndServe()

	tport := cfg.TcpPort

	taddr := ":" + strconv.Itoa(tport)

	log.Println("DNS Server Start at tcp", taddr)

	tcpServer = &dns.Server{Addr:taddr,Net:"tcp4",Handler:dnshandle}

	tcpServer.ListenAndServe()
}

func DNSServerStop()  {

	if udpServer != nil{
		udpServer.Shutdown()
		udpServer = nil
	}

	if tcpServer != nil{
		tcpServer.Shutdown()
		tcpServer = nil
	}

}
