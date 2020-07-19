package task

import pb "github.com/et-zone/gcelery/protos/base"

const (
	RES_TIMEOUT_ERR = "TimeOutError"
	UKNOWN_ERR      = "UnKnownError"
)

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

func NewResquest() *Request {
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
		} else {
			panic("SetKwargs err ,your kwargs not type as ['key','val']")
		}
	}
	request.Kwargs = kwgs
}

func (request *Request) SetReqBody(reqbody []byte) {
	request.ReqBody = reqbody
}

//内部返回响应信息
func GetPdResponse(res *pb.Response) Response {
	return Response{
		IsOk:    res.IsOk,
		Status:  res.Status,
		ResBody: res.ResBody,
	}

}

//获取响应信息
func SetResponseWithStatus(status string) Response {
	return Response{
		IsOk:    false,
		Status:  status,
		ResBody: []byte{},
	}

}

func TaskResponse(status string) *Response {
	return &Response{
		IsOk:    false,
		Status:  status,
		ResBody: []byte{},
	}

}
