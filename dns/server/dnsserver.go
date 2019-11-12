package server

import (
	"github.com/Ungigdu/BAS_contract_go/BAS_Ethereum"
	"github.com/kprc/basserver/config"
	"github.com/miekg/dns"
	"log"
	"net"
	"strconv"
)

var (
	dnshandle dns.HandlerFunc
)

func DnsHandle(writer dns.ResponseWriter, msg *dns.Msg) {
	//todo...
	m := msg.Copy()

	m.Compress = true
	m.Response = true

	A := &dns.A{}

	A.Hdr.Name = m.Question[0].Name
	A.Hdr.Rrtype = dns.TypeA
	A.Hdr.Class = dns.ClassINET
	A.Hdr.Ttl = 10
	A.Hdr.Rdlength = 4
	//A.A = net.ParseIP("123.56.153.221")

	//hash:=solsha3.SoliditySHA3(solsha3.String(A.Hdr.Name))
	qn := A.Hdr.Name

	if qn[len(qn)-1] == '.' {
		qn = qn[:len(qn)-1]

	}

	dr, err := BAS_Ethereum.QueryByString(qn)

	if err != nil {
		log.Println("Can't Get Domain Name info")
		m.Rcode = dns.RcodeBadKey

		writer.WriteMsg(m)

		return
	}

	A.A = net.IPv4(dr.IPv4[0], dr.IPv4[1], dr.IPv4[2], dr.IPv4[3])

	log.Println("Request Name: ", qn, A.A.String())
	log.Println("debug name: ", A.Hdr.Name)

	var rr []dns.RR

	rr = append(rr, A)

	m.Answer = rr

	writer.WriteMsg(m)

}

func DNSServerDaemon() {
	cfg := config.GetBasDCfg()
	uport := cfg.UpdPort
	uaddr := ":" + strconv.Itoa(uport)

	dnshandle = DnsHandle

	log.Println("DNS Server Start at udp", uaddr)

	go dns.ListenAndServe(uaddr, "udp4", dnshandle)

	tport := cfg.TcpPort

	taddr := ":" + strconv.Itoa(tport)

	log.Println("DNS Server Start at tcp", uaddr)

	log.Fatal(dns.ListenAndServe(taddr, "tcp4", dnshandle))
}
