package gcelery

import (
	"log"
	"net"

	// "time"

	"github.com/et-zone/gcelery/control"
	pb1 "github.com/et-zone/gcelery/protos/base"
	serv "github.com/et-zone/gcelery/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	TIMEOUT = 10
)

var gserver *GCeleryServer

type GCeleryServer struct {
	Server     *grpc.Server
	listen     net.Listener
	syncWroker *serv.SyncWroker
	cronWroker *serv.Cron
	// timeOut    int
}

func NewCelery(address string) *GCeleryServer {
	if gserver != nil {
		return gserver
	}
	server := &GCeleryServer{}
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	//grpc.ConnectionTimeout(time.Duration(5)) connectionTimeout,default=120s
	server.Server = grpc.NewServer()
	server.listen = listen
	// server.timeOut = TIMEOUT
	gserver = server
	return gserver
}

func NewTlsCelery(address string, cretFile string, key string) *GCeleryServer {
	creds, err := credentials.NewServerTLSFromFile(cretFile, key)
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}
	server := &GCeleryServer{}
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server.Server = grpc.NewServer(grpc.Creds(creds))
	server.listen = listen
	// server.timeOut = TIMEOUT
	return server
}

// func (this *GCeleryServer) SetTimeout(timeout int) {
// 	if timeout <= 0 {
// 		return
// 	}
// 	this.timeOut = timeout
// }

func (this *GCeleryServer) StartCelery() {
	if this.cronWroker != nil {
		serv.StartCron(this.cronWroker)
	}
	if this.syncWroker != nil {
		serv.StartSync(this.syncWroker)
	}
	if err := this.Server.Serve(this.listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//注册传输协议protobuf
func (this *GCeleryServer) RegisterTransport() {
	pb1.RegisterBridgeServer(this.Server, &control.GBase{})

}

func (this *GCeleryServer) RegisterCron(cronWroker *serv.Cron) {
	if this.cronWroker == nil {
		this.cronWroker = cronWroker
	}

}

func (this *GCeleryServer) RegisterSync(SyncWroker *serv.SyncWroker) {
	if this.syncWroker == nil {
		this.syncWroker = SyncWroker
	}
}

//grpc worker
func (this *GCeleryServer) InitCelery() {
	control.NewCeleryWorker()
}

func (this *GCeleryServer) RegisterCeleryWorker(fs ...func(*control.Request) (error, *control.Response)) {
	control.RegisterCeleryWorker(fs...)
}

//cron
func (this *GCeleryServer) NewCronWorker() *serv.Cron {
	return serv.NewCron()
}

//sync
func (this *GCeleryServer) NewSyncWroker() *serv.SyncWroker {
	return serv.InitSyncWroker()
}

//client Pool
func NewClientPool(max int, bindaddr string) *serv.CliPool {
	return serv.InitClientPool(max, bindaddr)
}

//client
func NewSTLClientPool(max int, bindaddr string, certFile string) *serv.CliPool {
	return serv.InitSTLClientPool(max, bindaddr, certFile)
}

//Request
func NewResQuest() *control.Request {
	return control.NewResQuest()
}
