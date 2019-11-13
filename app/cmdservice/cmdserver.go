package cmdservice

import (
	"google.golang.org/grpc"
	"sync"

	"net"

	"google.golang.org/grpc/reflection"
	"log"
	"errors"

	"github.com/kprc/basserver/config"
	"github.com/kprc/basserver/dns/server"
	"github.com/kprc/basserver/app/cmdpb"
	"github.com/kprc/basserver/app/cmdservice/api"
)

type cmdServer struct {
	localaddr  string
	grpcServer *grpc.Server
}

type CmdServerInter interface {
	StartCmdService()
	StopCmdService()
}

var (
	cmdServerInst     CmdServerInter
	cmdServerInstLock sync.Mutex
)

func GetCmdServerInst() CmdServerInter {
	if cmdServerInst == nil {
		cmdServerInstLock.Lock()
		defer cmdServerInstLock.Unlock()
		if cmdServerInst == nil {
			cfg := config.GetBasDCfg()
			cmdServerInst = &cmdServer{localaddr: cfg.CmdListenPort}
		}
	}

	return cmdServerInst
}

func (cs *cmdServer) checklocaladdress() error {
	if cs.localaddr == "" {
		return errors.New("No Server Listen address")
	}

	return nil
}

func (cs *cmdServer) StartCmdService() {
	if err := cs.checklocaladdress(); err != nil {
		log.Fatal("Start Cmd Service Failed", err)
		return
	}

	lis, err := net.Listen("tcp", cs.localaddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	cs.grpcServer = grpc.NewServer()

	cmdpb.RegisterDefaultcmdsrvServer(cs.grpcServer,&api.CmdDefaultServer{stop})

	reflection.Register(cs.grpcServer)
	log.Println("Commamd line server will start at", cs.localaddr)
	if err := cs.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (cs *cmdServer) StopCmdService() {
	config.GetBasDCfg().Save()
	server.DNSServerStop()

	cs.grpcServer.Stop()
	log.Println("Command line server stoped")
}

func stop() {
	GetCmdServerInst().StopCmdService()
}

