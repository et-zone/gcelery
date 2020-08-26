package worker

import (
	"fmt"

	// "golang.org/x/sys/unix"
	"github.com/et-zone/gcelery/task"
	// task "github.com/et-zone/gcelery/protos/base"
)

//base使用，传统的可带参数的业务线，随调随用，走rpc通信
func Do1(b *task.Request) (err error, r *task.Response) {

	// fmt.Println("拿到客户端传递的json: do1: ", string(b))
	// go func() {
	// 	select {
	// 	case <-ctx.Done():

	// 		err = errors.New("timeout")
	// 		byt = []byte{}
	// 		// panic("timeout")

	// 	}

	// }()

	// cout := 0
	// for {
	// 	cout++
	// 	log.Println("cout:", cout)
	// 	log.Println("method:", b.Method)
	// 	log.Println("Kwargs:", b.Kwargs)
	// 	log.Println("ReqBody:", string(b.ReqBody))
	// 	time.Sleep(time.Second)
	// 	if cout > 7 {
	// 		return nil, task.TaskResponse(task.UKNOWN_ERR)
	// 	}
	// }
	// log.Println("hhhhhh")
	// time.Sleep(time.Second)
	// time.Sleep(time.Second * 10)
	fmt.Println("do========task====", string(b.ReqBody))
	r = &task.Response{
		IsOk:    true,
		ResBody: b.ReqBody,
		Status:  "ok",
	}

	return nil, r
}

// func Do2(b []byte) (error, []byte) {
// 	fmt.Println("拿到客户端传递的json: do2: ", string(b))
// 	return nil, b
// }
// func Do3(b []byte) (error, []byte) {
// 	fmt.Println("拿到客户端传递的json: do3: ", string(b))
// 	return nil, b
// }
