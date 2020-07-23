package server

import (
	"strings"
	"time"

	"context"
	"log"

	pb "github.com/et-zone/gcelery/protos/base"
	"github.com/et-zone/gcelery/task"
	// "google.golang.org/grpc"
)

// func NewCurSor(addr string) Cursor {
// 	c := &cursor{
// 		conn: newConn(addr),
// 	}
// 	c.cursor = pb.NewBridgeClient(c.conn)
// 	c.timeout = 5
// 	return c
// }

type Cursor interface {
	Do(req *task.Request) task.Response
	doContext(ctx context.Context, req *task.Request) task.Response
}

//cursor
type cursor struct {
	cursor pb.BridgeClient
	// conn    *grpc.ClientConn
	timeout int
}

func (this *cursor) Do(req *task.Request) task.Response {
	if this == nil {
		log.Fatal("Do err, Cursor is nil can not Do function ")
	}
	t := time.Now()
	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(this.timeout))
	r, err := this.cursor.Dao(ctx, &pb.Request{ //context.TODO()==default
		Method:  req.Method,
		ReqBody: req.ReqBody,
		Kwargs:  req.Kwargs,
	})
	log.Println(time.Since(t))
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
func (this *cursor) doContext(ctx context.Context, req *task.Request) task.Response {
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

// func (this *cursor) Close() {
// 	client.pool.Lock()
// 	client.pool.Pool.Put(this.conn)
// 	this.conn = nil
// 	this = nil
// }
