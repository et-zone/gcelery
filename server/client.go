package server

import (
	"log"
	"sync"

	pb "github.com/et-zone/gcelery/protos/base"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var cliPool *CliPool

type CeleryClient struct {
	conn *grpc.ClientConn
}

type CliPool struct {
	sync.Pool
	Maxconn  int //连接数
	Address  string
	isSTL    bool
	certFile string //cert文件路径
}

type Msg struct {
	Status  string `json:"status" bson:"status"`
	Message []byte `json:"message" bson:"message"`
}

func NewClient(address string) *CeleryClient {
	con := &CeleryClient{}
	con.conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	return con
}
func NewTlsClient(address string, certFile string) *CeleryClient {
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		log.Fatalf("could not process the credentials: %v", err)
	}
	con := &CeleryClient{}
	con.conn, err = grpc.Dial(address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatal(err.Error())
	}
	return con
}

//初始化client连接池
func InitClientPool(max int, addr string) *CliPool {
	if max <= 0 {
		panic("initClientPool err, Maxconn can not < 0")
	}
	if cliPool != nil {
		return cliPool
	}
	cliPool = &CliPool{Pool: sync.Pool{
		New: func() interface{} {
			return nil
		}},
		Maxconn: max,
		Address: addr,
		isSTL:   false,
	}
	for i := 0; i < cliPool.Maxconn; i++ {
		cliPool.Put(func() interface{} {
			return NewClient(cliPool.Address)
		}())
	}

	return cliPool
}

//初始化STL client连接池
func InitSTLClientPool(max int, addr string, certFile string) *CliPool {
	if max <= 0 {
		panic("initClientPool err, Maxconn can not < 0")
	}
	if cliPool != nil {
		return cliPool
	}
	cliPool = &CliPool{Pool: sync.Pool{
		New: func() interface{} {
			return nil
		}},
		Maxconn:  max,
		Address:  addr,
		isSTL:    true,
		certFile: certFile,
	}

	for i := 0; i < cliPool.Maxconn; i++ {
		cliPool.Put(func() interface{} {
			return NewTlsClient(cliPool.Address, certFile)
		}())
	}

	return cliPool
}

func (clipool *CliPool) GetCursor() *Cursor {
	client, isok := clipool.Get().(*CeleryClient)
	if !isok { //新增一个
		if clipool.isSTL == true {
			client := NewTlsClient(cliPool.Address, clipool.certFile)
			cli := pb.NewBridgeClient(client.conn)
			return &Cursor{cli, client}
		} else {
			client := NewClient(cliPool.Address)
			cli := pb.NewBridgeClient(client.conn)
			return &Cursor{cli, client}
		}

	}
	cli := pb.NewBridgeClient(client.conn)

	return &Cursor{cli, client}
}

func (clipool *CliPool) ClosePool() {
	for {
		client, isok := clipool.Get().(*CeleryClient)
		if !isok {
			return
		}
		err := client.conn.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}

}
