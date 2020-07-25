package test

import (
	// "fmt"
	// "strings"
	// "time"

	"github.com/et-zone/gcelery/server"
	"github.com/et-zone/gcelery/task"

	"testing"
	// "context"
	// "time"

	"github.com/et-zone/gcelery"
)

const defaultMethod = "Do1"

var pool *server.CeleryClient
var req *task.Request = &task.Request{
	Method:  defaultMethod,
	ReqBody: []byte("你好"),
	Kwargs:  map[string]string{"aaa": "111", "bbb": "222"},
}

func Doin() {
	// count := 0

	// t := time.Now().Add(time.Second * time.Duration(10))
	cons := pool.Clone()
	cons.Do(req)
	// if strings.Contains(r.Status, "connection refused") {
	// 	fmt.Println("connection refused")
	// 	break
	// }
	// if time.Now().After(t) {
	// 	fmt.Println(count)
	// 	break
	// }
	// cons.Close()

}

//go test t_test.go -bench=.

//连接池+普通版本
func Benchmark(b *testing.B) {
	addr := "localhost:50051"
	// addr := "49.232.190.114:50051"
	// addr := "10.206.0.15:50051"
	pool = gcelery.NewClient(addr)
	defer pool.Close()
	for i := 0; i < b.N; i++ { // b.N，测试循环次数
		Doin()
	}
	// pool.SetTimeOut(5)

	// for i := 0; i < 20; i++ {
	// 	go func() {
	// 		count := 0
	// 		cursor := pool.Clone()
	// 		t := time.Now().Add(time.Second * time.Duration(10))
	// 		for {
	// 			count++
	// 			r := cursor.Do(req)
	// 			if strings.Contains(r.Status, "connection refused") {
	// 				fmt.Println("connection refused")
	// 				break
	// 			}
	// 			if time.Now().After(t) {
	// 				fmt.Println(count)
	// 				break
	// 			}
	// 		}

	// 		cursor.Close()
	// 		// fmt.Println(time.Since(t))
	// 	}()

	// }

	// msg := cursor.Do(req)

	// fmt.Println(msg.Status, string(msg.ResBody))

	// time.Sleep(time.Second * 15)
}
