package server

import (
	"context"
	"log"

	pb "github.com/et-zone/gcelery/protos/base"
)

var err error

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
	return msg
}

func (this *Cursor) DoContext(ctx context.Context, method string, data []byte, args ...string) Msg {
	if this == nil {
		log.Fatal("Do err, Cursor is nil can not Do function ")
	}
	r, err := this.cursor.Dao(ctx, &pb.Request{ //context.TODO()
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
