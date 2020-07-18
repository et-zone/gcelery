package control

import (
	"context"
	// "fmt"
	"log"

	// "time"

	pb "github.com/et-zone/gcelery/protos/base"
	// wk "github.com/et-zone/gcelery/server"
)

type GBase struct{}

const (
	RES_TIMEOUT_ERR = "TimeOutError"
	UKNOWN_ERR      = "UnKnownError"
)

func (s *GBase) Dao(ctx context.Context, in *pb.Request) (*pb.Response, error) {
	// log.Printf("Received1: %v", in.Method)
	// ctx, _ = context.WithTimeout(context.Background(), time.Duration(5)*time.Second)

	err, b := s.BaseFunc(in, GetCeleryWorker(in.Method))
	if err != nil {
		return &pb.Response{Status: "false", ResBody: b.ResBody}, err
	}
	// return &pb.Response{Status: "111===" + in.Method + " b:" + string(b)}, nil
	return &pb.Response{Status: "success", ResBody: b.ResBody}, nil

}

// 不带参数
// func (s *GBase) BaseFunc(f func() error) error {
// 	return f()
// }

//可以带参数json
func (s *GBase) BaseFunc(r *pb.Request, f func(*Request) (error, *Response)) (error, *Response) {
	if f == nil {
		log.Fatal("BaseFunc err, variable f eq nil")
	}
	req := &Request{
		Method:  r.Method,
		Kwargs:  r.Kwargs,
		ReqBody: r.ReqBody,
	}
	return f(req)
}

type Request struct {
	Method  string            `bson:"method" json:"method"`
	Kwargs  map[string]string `bson:"kwargs" json:"kwargs"`
	ReqBody []byte            `bson:"reqBody" json:"reqBody"`
}
type Response struct {
	IsOk    bool   `bson:"isOk" json:"isOk"`
	Status  string `bson:"status" json:"status"`
	ResBody []byte `bson:"resBody" json:"resBody"`
}

func NewResQuest() *Request {
	return &Request{}
}

func (request *Request) SetMethod(method string) {
	request.Method = method
}

func (request *Request) SetKwargs(kwargs ...[]string) {
	kwgs := map[string]string{}
	for _, args := range kwargs {
		if len(args) == 2 {
			kwgs[args[0]] = args[1]
		}
	}
	request.Kwargs = kwgs
}

func (request *Request) SetReqBody(reqbody []byte) {
	request.ReqBody = reqbody
}

//根据返回响应返回响应信息
func GetResponse(res *pb.Response) Response {
	return Response{
		IsOk:    res.IsOk,
		Status:  res.Status,
		ResBody: res.ResBody,
	}

}

func GetErrResponse(errStatus string) Response {
	return Response{
		IsOk:    false,
		Status:  errStatus,
		ResBody: []byte{},
	}

}

func SetWResponse(Status string) *Response {
	return &Response{
		IsOk:    false,
		Status:  Status,
		ResBody: []byte{},
	}

}
