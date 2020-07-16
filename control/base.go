package control

import (
	"context"
	// "fmt"
	"log"

	pb "github.com/et-zone/gcelery/protos/base"
	wk "github.com/et-zone/gcelery/server"
)

type GBase struct{}

func (s *GBase) Dao(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	// log.Printf("Received1: %v", in.Method)

	err, b := s.BaseFunc(ctx, in.Data, wk.GetRpcWorker(in.Method))
	if err != nil {
		return &pb.Response{Status: "false", Msg: b}, err
	}
	// return &pb.Response{Status: "111===" + in.Method + " b:" + string(b)}, nil
	return &pb.Response{Status: "success", Msg: b}, nil

}

// 不带参数
// func (s *GBase) BaseFunc(f func() error) error {
// 	return f()
// }

//可以带参数json
func (s *GBase) BaseFunc(ctx context.Context, b []byte, f func(context.Context, []byte) (error, []byte)) (error, []byte) {
	if f == nil {
		log.Fatal("BaseFunc err, variable f eq nil")
	}
	return f(ctx, b)
}
