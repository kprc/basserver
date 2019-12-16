package server

import (
	"github.com/Ungigdu/BAS_contract_go/BAS_Ethereum"
	"github.com/kprc/basserver/config"

	"encoding/binary"
	"github.com/btcsuite/btcutil/base58"
	"github.com/miekg/dns"
	"github.com/pkg/errors"
	"log"
	"net"
	"strconv"
)

const (
	TypeBCAddr = 65
)

var (
	dnshandle dns.HandlerFunc
)

var (
	udpServer *dns.Server
	tcpServer *dns.Server
)

type DR struct {
	BAS_Ethereum.DomainRecord
}

func (dr *DR) IntIPv4() uint32 {
	return binary.BigEndian.Uint32(dr.IPv4[:])
}


func sendErrMsg(w dns.ResponseWriter, msg *dns.Msg, errCode int) {
	m := msg.Copy()

	m.Compress = true
	m.Response = true

	m.Rcode = errCode

	w.WriteMsg(m)
}

func buildAnswer(ipv4 [4]byte, q dns.Question) []dns.RR {
	A := &dns.A{}

	A.Hdr.Name = q.Name
	A.Hdr.Rrtype = dns.TypeA
	A.Hdr.Class = dns.ClassINET
	A.Hdr.Ttl = 10
	A.Hdr.Rdlength = 4

	A.A = net.IPv4(ipv4[0], ipv4[1], ipv4[2], ipv4[3])

	log.Println("Request Name: ", q.Name, A.A.String())

	var rr []dns.RR

	rr = append(rr, A)

	return rr
}

func replyTypA(w dns.ResponseWriter, msg *dns.Msg, q dns.Question) error {

	qn := q.Name
	if qn[len(qn)-1] == '.' {
		qn = qn[:len(qn)-1]

	}
	if bdr, err := BAS_Ethereum.QueryByString(qn); err != nil {
		return errors.New("Query DN from Ethereum error: " + err.Error())
	} else {
		dr := &DR{bdr}
		if dr.IntIPv4() == 0 {
			return errors.New("ipv4 address error")
		}
		m := msg.Copy()
		m.Compress = true
		m.Response = true

		m.Answer = buildAnswer(dr.IPv4, q)

		w.WriteMsg(m)

		return nil
	}
}

func replyTraditionTypA(w dns.ResponseWriter, msg *dns.Msg) {
	cnt := 0

	for {

		cnt++

		s := GetDns()

		if s == "" {
			sendErrMsg(w, msg, dns.RcodeServerFailure)
			return
		}
		//log.Println(s+":53",msg.Question[0].Name)
		if m, err := dns.Exchange(msg, s+":53"); err != nil {
			FailDns(s)
			log.Println("failed "+s+":53", msg.Question[0].Name)
			if cnt >= MaxTimes() {
				sendErrMsg(w, msg, dns.RcodeBadKey)
				return
			}
		} else {
			w.WriteMsg(m)
			log.Println("success "+s+":53", msg.Question[0].Name)
			return
		}

	}
}

func replyTypPTR(w dns.ResponseWriter, msg *dns.Msg, q dns.Question) error {
	return nil
}

func replyTraditionTypPTR(w dns.ResponseWriter, msg *dns.Msg, q dns.Question) {

}

func replyTypBCA(w dns.ResponseWriter, msg *dns.Msg, q dns.Question) error {
	qn := q.Name
	if qn[len(qn)-1] == '.' {
		qn = qn[:len(qn)-1]

	}

	var b []byte
	b = base58.Decode(qn)
	var barr [32]byte

	for i := 0; i < len(b); i++ {
		barr[i] = b[i]
	}

	if bdr, err := BAS_Ethereum.QueryByBCAddress(barr); err != nil {
		return errors.New("Query BCA from Ethereum error: " + err.Error())
	} else {
		dr := &DR{bdr}
		if dr.IntIPv4() == 0 {
			return errors.New("ipv4 address error")
		}
		m := msg.Copy()
		m.Compress = true
		m.Response = true

		m.Answer = buildAnswer(dr.IPv4, q)

		w.WriteMsg(m)

		return nil
	}
}

func DnsHandleTradition(w dns.ResponseWriter, msg *dns.Msg) {
	if len(msg.Question) == 0 {
		sendErrMsg(w, msg, dns.RcodeFormatError)
		return
	}
	q := msg.Question[0]

	if q.Qclass != dns.ClassINET {
		sendErrMsg(w, msg, dns.RcodeNotImplemented)
		return
	}

	switch q.Qtype {
	case dns.TypeA:
		if err := replyTypA(w, msg, q); err != nil {
			replyTraditionTypA(w, msg)
		}
	case dns.TypePTR:
		//if err:=replyTypPTR(w,msg,q);err!=nil{
		//	replyTraditionTypPTR(w,msg,q)
		//}
		replyTraditionTypA(w, msg)
	case TypeBCAddr:

		if err := replyTypBCA(w, msg, q); err != nil {
			//replyTraditionTypPTR(w,msg,q)
			sendErrMsg(w, msg, dns.RcodeBadKey)
		}

	default:
		//sendErrMsg(w,msg,dns.RcodeNotImplemented)
		replyTraditionTypA(w, msg)
		return
	}

}

func DNSServerDaemon() {
	cfg := config.GetBasDCfg()

	//fmt.Println(*cfg)

	uport := cfg.UpdPort
	uaddr := ":" + strconv.Itoa(uport)

	dnshandle = DnsHandleTradition

	log.Println("DNS Server Start at udp", uaddr)

	udpServer = &dns.Server{}
	udpServer.Addr = uaddr
	udpServer.Handler = dnshandle
	udpServer.Net = "udp4"

	go udpServer.ListenAndServe()

	tport := cfg.TcpPort

	taddr := ":" + strconv.Itoa(tport)

	log.Println("DNS Server Start at tcp", taddr)

	tcpServer = &dns.Server{Addr: taddr, Net: "tcp4", Handler: dnshandle}

	tcpServer.ListenAndServe()
}

func DNSServerStop() {

	if udpServer != nil {
		udpServer.Shutdown()
		udpServer = nil
	}

	if tcpServer != nil {
		tcpServer.Shutdown()
		tcpServer = nil
	}

}
