package server

import (
	"log"
	"strings"
	"sync"

	"time"

	"context"

	// "github.com/et-zone/gcelery/control"
	pb "github.com/et-zone/gcelery/protos/base"
	"github.com/et-zone/gcelery/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var cliPool *CliPool
var err error

const TimeOut = 60

type CeleryClient struct {
	conn    *grpc.ClientConn
	timeout int
	isnew   bool
}

type CliPool struct {
	sync.Mutex
	sync.Pool
	Maxconn     int    //连接数
	Address     string //
	isSTL       bool   //stl
	certFile    string //cert文件路径
	timeout     int    //超时时间
	MaxIdleConn int    //最大连接数
	MinOpenConn int    //最小存活数
	LocalConn   int    //当前连接数

}

func NewClient(address string) *CeleryClient {
	con := &CeleryClient{}
	//, grpc.WithTimeout(time.Duration(10)*time.Second) client
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

//init client Pool
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
		Mutex:   sync.Mutex{},
		Maxconn: max,
		Address: addr,
		isSTL:   false,
		timeout: TimeOut,
	}

	for i := 0; i < cliPool.Maxconn; i++ {
		cliPool.Put(func() interface{} {
			return NewClient(cliPool.Address)
		}())
	}

	return cliPool
}

//init STL client Pool
func InitSTLClientPool(max int, addr string, certFile string) *CliPool {
	if max <= 0 {

	}
	if max > 50 {
		panic("initClientPool err, Maxconn can not < 0")
	}
	if cliPool != nil {
		return cliPool
	}
	cliPool = &CliPool{Pool: sync.Pool{
		New: func() interface{} {
			return nil
		}},
		Mutex:    sync.Mutex{},
		Maxconn:  max,
		Address:  addr,
		isSTL:    true,
		certFile: certFile,
		timeout:  TimeOut,
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
			client.timeout = clipool.timeout
			cli := pb.NewBridgeClient(client.conn)
			return &Cursor{cli, client, true}
		} else {
			client := NewClient(cliPool.Address)
			client.timeout = clipool.timeout
			cli := pb.NewBridgeClient(client.conn)
			return &Cursor{cli, client, true}
		}

	}
	client.timeout = clipool.timeout
	cli := pb.NewBridgeClient(client.conn)

	return &Cursor{cli, client, true}
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

//SetClient Timeout default 60s
func (clipool *CliPool) SetTimeOut(timeout int) {
	if timeout <= 0 {
		log.Println("set timeout ,timeout can not <= 0")
	}
	cliPool.timeout = timeout
}

//cursor
type Cursor struct {
	cursor       pb.BridgeClient
	client       *CeleryClient
	isPoolCursor bool
}

func (this *Cursor) Do(req *task.Request) task.Response {
	if this == nil {
		log.Fatal("Do err, Cursor is nil can not Do function ")
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(this.client.timeout))
	r, err := this.cursor.Dao(ctx, &pb.Request{ //context.TODO()==default
		Method:  req.Method,
		ReqBody: req.ReqBody,
		Kwargs:  req.Kwargs,
	})
	if err != nil {
		// log.Println(err.Error())
		if strings.Contains(err.Error(), "DeadlineExceeded") {
			return task.SetResponseWithStatus(task.RES_TIMEOUT_ERR)
		} else {
			return task.SetResponseWithStatus(err.Error())
		}
	}

	return task.GetPdResponse(r)
}

/*
exampel:
	自定义控制程context 超时时间
	// ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(5))
	// ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Second*5))
	// ctx, cancel := context.WithCancel(context.Background()) //程序退出才会走下一步
	// defer cancel()
	// msg := cursor.DoContext(ctx, req)

*/
//自定义控制程context 暂时没必要，因为已经可以设置自定义超时时间了，暂且环比该函数
func (this *Cursor) doContext(ctx context.Context, req *task.Request) task.Response {
	if this == nil {
		log.Fatal("Do err, Cursor is nil can not Do function ")
	}

	r, err := this.cursor.Dao(ctx, &pb.Request{
		Method:  req.Method,
		ReqBody: req.ReqBody,
		Kwargs:  req.Kwargs,
	})
	if err != nil {
		// log.Println(err.Error())
		if strings.Contains(err.Error(), "DeadlineExceeded") {
			return task.SetResponseWithStatus(task.RES_TIMEOUT_ERR)
		} else {
			return task.SetResponseWithStatus(err.Error())
		}

	}

	return task.GetPdResponse(r)
}

func (this *Cursor) Close() {

	if this.isPoolCursor == true {
		if cliPool == nil {
			log.Println("Close Cursor err , not found CliPool")
			return
		}
		cliPool.Put(this.client)
	}
	this = nil
}

func (celeryClient *CeleryClient) Close() {
	if celeryClient == nil {
		log.Println("client err  client==nil can not Close")
		return
	}
	err := celeryClient.conn.Close()
	if err != nil {
		log.Println(err.Error())
	}
}

func (celeryClient *CeleryClient) GetCursor() *Cursor {
	if celeryClient == nil {
		log.Println("client err  client==nil can not get Cursor")
		return nil
	}
	celeryClient.timeout = TimeOut //default
	cli := pb.NewBridgeClient(celeryClient.conn)

	return &Cursor{cli, celeryClient, false}
}

//SetClient Timeout default 60s
func (celeryClient *CeleryClient) SetTimeOut(timeout int) {
	if timeout <= 0 {
		log.Println("set timeout ,timeout can not <= 0")
	}
	celeryClient.timeout = timeout
}

//client init potion,  not used
type CliOption struct {
	name  string
	value interface{}
}

func NewCliOpt(argName string, val interface{}) CliOption {
	return CliOption{
		name:  argName,
		value: val,
	}
}
