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

var gserver *GrpcServer

//测试动态价值对应配置
type GrpcServer struct {
	Server *grpc.Server
	listen net.Listener
}

func NewGrpc(address string) *GrpcServer {
	if gserver != nil {
		return gserver
	}
	server := &GrpcServer{}
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server.Server = grpc.NewServer()
	server.listen = listen
	gserver = server
	return gserver
}
func NewTlsGrpc(address string, cretFile string, key string) *GrpcServer {
	creds, err := credentials.NewServerTLSFromFile(cretFile, key)
	if err != nil {
		log.Fatalf("Failed to setup TLS: %v", err)
	}
	server := &GrpcServer{}
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server.Server = grpc.NewServer(grpc.Creds(creds))
	server.listen = listen
	return server
}

func (this *GrpcServer) StartGrpc() {
	if err := this.Server.Serve(this.listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (this *GrpcServer) RegisterBridgeBaseServer() {
	pb1.RegisterBridgeServer(this.Server, &control.GBase{})
}

//grpc worker
func (this *GrpcServer) NewRpcWorker() {
	serv.NewRpcWorker()
}

func (this *GrpcServer) RegisterRpcWorker(fs ...func([]byte) (error, []byte)) {
	serv.RegisterRpcWorker(fs...)
}

//cron
func (this *GrpcServer) NewCronWorker() *serv.Cron {
	return serv.NewCron()
}

//sync
func (this *GrpcServer) NewSyncWroker() *serv.SyncWroker {
	return serv.InitSyncWroker()
}

//client
func NewClient(bindaddr string) *serv.GConnect {
	return serv.NewClient(bindaddr)
}

//client
func NewTlsClient(bindaddr string, certFile string) *serv.GConnect {
	return serv.NewTlsClient(bindaddr, certFile)
}
