package server

import (
	"context"
	"log"

	pb "github.com/et-zone/gcelery/protos/base"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials"
)

var err error

// type CeleryClient struct {
// 	conn *grpc.ClientConn
// }
// type Msg struct {
// 	Status  string `json:"status" bson:"status"`
// 	Message []byte `json:"message" bson:"message"`
// }

// func NewClient(address string) *CeleryClient {
// 	con := &CeleryClient{}
// 	con.conn, err = grpc.Dial(address, grpc.WithInsecure())
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	return con
// }
// func NewTlsClient(address string, certFile string) *CeleryClient {
// 	creds, err := credentials.NewClientTLSFromFile(certFile, "")
// 	if err != nil {
// 		log.Fatalf("could not process the credentials: %v", err)
// 	}
// 	con := &CeleryClient{}
// 	con.conn, err = grpc.Dial(address, grpc.WithTransportCredentials(creds))
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}
// 	return con
// }

// func (this *CeleryClient) CloseClient() {
// 	this.conn.Close()
// }

// func (this *CeleryClient) GetCursor() *Cursor {
// 	cli := pb.NewBridgeClient(this.conn)

// 	return &Cursor{cli}
// }

//cursor
type Cursor struct {
	cursor pb.BridgeClient
	client *CeleryClient
}

func (this *Cursor) Do(method string, data []byte, args ...string) Msg {
	if this == nil {
		log.Fatal("Do err, Cursor is nil can not Do function ")
	}
	r, err := this.cursor.Dao(context.TODO(), &pb.Request{ //context.TODO()
		Method: method,
		Data:   data,
		Args:   args,
	})
	if err != nil {
		log.Println(err.Error())
	}
	msg := Msg{
		Status:  r.Status,
		Message: r.Msg,
	}
	// log.Printf("Greeting: %s", r.Status)
	// time.Sleep(time.Second)
	return msg
}
func (this *Cursor) Close() {
	if cliPool == nil {
		log.Println("Close Cursor err , not found CliPool")
		return
	}
	cliPool.Put(this.client)
	this = nil
}
