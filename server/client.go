package server

import (
	"log"
	"sync"

	// "github.com/et-zone/gcelery/control"
	pb "github.com/et-zone/gcelery/protos/base"
	"github.com/et-zone/gcelery/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var client *CeleryClient
var err error

const TimeOut = 60

type CeleryClient struct {
	conn    *grpc.ClientConn
	timeout int
	isnew   bool
	pool    *CliPool
	// Cursor
}

type CliPool struct {
	sync.Mutex
	conn        *grpc.ClientConn
	Maxconn     int    //连接数
	Address     string //
	isSTL       bool   //stl
	certFile    string //cert文件路径
	cred        *credentials.TransportCredentials
	timeout     int //超时时间
	MaxIdleConn int //最大连接数
	MinOpenConn int //最小存活数
	LocalConn   int //当前连接数

}

func newConn(address string) *grpc.ClientConn {

	//, grpc.WithTimeout(time.Duration(10)*time.Second) client
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	return conn
}
func newTlsConn(address string, certFile string) (*grpc.ClientConn, *credentials.TransportCredentials) {
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err.Error())
	}
	return conn, &creds
}

//init client Pool
func InitClient(addr string) *CeleryClient {

	if client != nil {
		return client
	}
	client = &CeleryClient{
		// Cursor: nil,
	}
	cliPool := &CliPool{
		conn:    newConn(addr),
		Mutex:   sync.Mutex{},
		Address: addr,
		isSTL:   false,
		timeout: TimeOut,
	}
	client.pool = cliPool
	client.timeout = cliPool.timeout
	return client
}

//init STL client Pool
func InitSTLClient(addr string, certFile string) *CeleryClient {
	if client != nil {
		return client
	}
	client = &CeleryClient{ /*Cursor: nil*/ }
	conn, cred := newTlsConn(addr, certFile)
	cliPool := &CliPool{
		conn:     conn,
		cred:     cred,
		Mutex:    sync.Mutex{},
		Address:  addr,
		isSTL:    true,
		certFile: certFile,
		timeout:  TimeOut,
	}
	client.pool = cliPool
	client.timeout = cliPool.timeout

	return client
}

func (cli *CeleryClient) Close() {
	var err error
	cli.conn = nil
	// cli.Cursor = nil
	if cli.pool != nil {
		err = cli.pool.conn.Close()
	}
	if err != nil {
		log.Println(err.Error())
	}

}

func (cli *CeleryClient) cursor() Cursor {
	return &cursor{pb.NewBridgeClient(client.pool.conn), client.pool.timeout}
}

func (cli *CeleryClient) Do(req *task.Request) task.Response {
	return cli.cursor().Do(req)
}
