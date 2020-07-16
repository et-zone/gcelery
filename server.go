package gcelery

import (
	"log"
	"net"

	"github.com/et-zone/gcelery/control"
	pb1 "github.com/et-zone/gcelery/protos/base"
	serv "github.com/et-zone/gcelery/server"
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

func NewCelery(address string) *GCeleryServer {
	if gserver != nil {
		return gserver
	}
	server := &GCeleryServer{}
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server.Server = grpc.NewServer()
	server.listen = listen
	gserver = server
	return gserver
}
func NewTlsGrpc(address string, cretFile string, key string) *GCeleryServer {
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

func (this *GCeleryServer) StartCelery() {
	if this.cronWroker != nil {
		this.cronWroker.Start()
	}
	if this.syncWroker != nil {
		this.syncWroker.Start()
	}

	if err := this.Server.Serve(this.listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//注册传输协议protobuf
func (this *GCeleryServer) RegisterTransport() {
	pb1.RegisterTransport(this.Server, &control.GBase{})
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
	serv.NewCeleryWorker()
}

func (this *GCeleryServer) RegisterCeleryWorker(fs ...func([]byte) (error, []byte)) {
	serv.RegisterRpcWorker(fs...)
}

//cron
func (this *GCeleryServer) NewCronWorker() *serv.Cron {
	return serv.NewCron()
}

//sync
func (this *GCeleryServer) NewSyncWroker() *serv.SyncWroker {
	return serv.InitSyncWroker()
}

//client
func NewClient(bindaddr string) *serv.CeleryClient {
	return serv.NewClient(bindaddr)
}

//client
func NewTlsClient(bindaddr string, certFile string) *serv.CeleryClient {
	return serv.NewTlsClient(bindaddr, certFile)
}
