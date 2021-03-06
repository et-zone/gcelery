package control

import (
	"context"
	"log"

	pb "github.com/et-zone/gcelery/protos/base"
	"github.com/et-zone/gcelery/task"
)

type GBase struct {
	// pb.UnimplementedBridgeServer
}

func (s *GBase) Dao(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	// log.Printf("Received1: %v", in.Method)
	// ctx, _ = context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	// t := time.Now()
	err, b := s.BaseFunc(in, GetCeleryWorker(in.Method))
	// log.Println(time.Since(t))
	if err != nil {
		return &pb.Response{IsOk: false, Status: "req fail", ResBody: b.ResBody}, err
	}
	// return &pb.Response{Status: "111===" + in.Method + " b:" + string(b)}, nil
	return &pb.Response{IsOk: true, Status: "req success", ResBody: b.ResBody}, nil

}

// Without Args
// func (s *GBase) BaseFunc(f func() error) error {
// 	return f()
// }

// Have  Args
func (s *GBase) BaseFunc(r *pb.Request, f func(*task.Request) (error, *task.Response)) (error, *task.Response) {
	if f == nil {
		log.Fatal("BaseFunc err, variable f eq nil")
	}

	req := &task.Request{
		Method:  r.Method,
		Kwargs:  r.Kwargs,
		ReqBody: r.ReqBody,
	}
	return f(req)
}
