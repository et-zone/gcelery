package gcelery

import (
	"log"
	"net"

	// "time"

	"github.com/et-zone/gcelery/control"
	pb1 "github.com/et-zone/gcelery/protos/base"
	serv "github.com/et-zone/gcelery/server"
	"github.com/et-zone/gcelery/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var gserver *GCeleryServer

type GCeleryServer struct {
	Server     *grpc.Server
	listen     net.Listener
	syncWroker *serv.SyncWroker
	cronWroker *serv.Cron
}

//New Server
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
	gserver = server
	return gserver
}

//New TLS Server
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
	return server
}

//Start Server GCelery
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

//Register Protobuf
func (this *GCeleryServer) RegisterTransport() {
	pb1.RegisterBridgeServer(this.Server, &control.GBase{})

}

// Register Cron Wroker to Server
func (this *GCeleryServer) RegisterCron(cronWroker *serv.Cron) {
	if this.cronWroker == nil {
		this.cronWroker = cronWroker
	}

}

//Register Sync Long Wroker to Server
func (this *GCeleryServer) RegisterSync(SyncWroker *serv.SyncWroker) {
	if this.syncWroker == nil {
		this.syncWroker = SyncWroker
	}
}

//Gcelery Task Wroker init
func (this *GCeleryServer) InitCelery() {
	control.NewCeleryWorker()
}

func (this *GCeleryServer) RegisterCeleryWorker(fs ...func(*task.Request) (error, *task.Response)) {
	control.RegisterCeleryWorker(fs...)
}

//Cron Wroker
func (this *GCeleryServer) NewCronWorker() *serv.Cron {
	return serv.NewCron()
}

//Sync Long wroker
func (this *GCeleryServer) NewSyncWroker() *serv.SyncWroker {
	return serv.InitSyncWroker()
}

//Client
func NewClient(bindaddr string) *serv.CeleryClient {
	return serv.InitClient(bindaddr)
}

//STL Client
func NewSTLClient(bindaddr string, certFile string) *serv.CeleryClient {
	return serv.InitSTLClient(bindaddr, certFile)
}

//Request
func NewTaskResquest() *task.Request {
	return task.NewResquest()
}
