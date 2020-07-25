package main

import (
	"fmt"
	"strings"
	"time"

	// "testing"
	// "context"
	// "time"

	"github.com/et-zone/gcelery"
)

//连接池+普通版本
func TestPool1() {
	addr := "localhost:50051"
	// addr := "10.206.0.15:50051"
	// addr := "49.232.190.114:50051"
	defaultMethod := "Do1"
	pool := gcelery.NewClient(addr)
	defer pool.Close()

	req := gcelery.NewTaskResquest()
	req.Method = defaultMethod
	req.ReqBody = []byte("你好")
	req.Kwargs = map[string]string{
		"aaa": "111",
		"bbb": "222",
	}
	for i := 0; i < 10; i++ {
		go func() {
			count := 0
			client := pool.Clone()
			t := time.Now().Add(time.Second * time.Duration(10))
			for {
				count++
				r := client.Do(req)
				if strings.Contains(r.Status, "connection refused") {
					fmt.Println("connection refused")
					break
				}
				if time.Now().After(t) {
					fmt.Println(count)
					break
				}
			}

			// fmt.Println(time.Since(t))
		}()

	}

	// msg := cursor.Do(req)

	// fmt.Println(msg.Status, string(msg.ResBody))

	time.Sleep(time.Second * 15)
}

//连接池+STL 安全
func TestPool2() {
	addr := "gzytobe.cn:50051"
	defaultMethod := "Do1"
	pool := gcelery.NewSTLClient(addr, "/tl/server.crt")
	defer pool.Close()
	client := pool.Clone()
	req := gcelery.NewTaskResquest()
	req.Method = defaultMethod
	req.ReqBody = []byte("你好")
	req.Kwargs = map[string]string{
		"aaa": "111",
		"bbb": "222",
	}
	cout := 0
	t := time.Now().Add(time.Second * 1)
	for { //单秒4200左右

		cout++
		/*msg :=*/ client.Do(req)
		// fmt.Println(msg.Status, string(msg.ResBody))
		// if cout > 10000 {
		// 	break
		// }
		if time.Now().After(t) {
			fmt.Println("cout:", cout)
			return
		}
	}

}

//普通单连接,单连接太多了影响性能，量不能太多
func TestClient1() {

	addr := "localhost:50051"
	defaultMethod := "Do1"
	pool := gcelery.NewClient(addr)
	defer pool.Close()
	client := pool.Clone()
	req := gcelery.NewTaskResquest()
	req.Method = defaultMethod
	req.ReqBody = []byte("你好")
	req.Kwargs = map[string]string{
		"aaa": "111",
		"bbb": "222",
	}
	msg := client.Do(req)

	fmt.Println(msg.Status, string(msg.ResBody))
}

//STL 单连接
func TestClient2() {
	addr := "gzytobe.cn:50051"
	defaultMethod := "Do1"
	pool := gcelery.NewSTLClient(addr, "/tl/server.crt")
	defer pool.Close()
	client := pool.Clone()
	//==================================================
	req := gcelery.NewTaskResquest()
	req.Method = defaultMethod
	req.ReqBody = []byte("你好")
	req.Kwargs = map[string]string{
		"aaa": "111",
		"bbb": "222",
	}
	cout := 0
	t := time.Now().Add(time.Second)
	for { //单秒4200左右
		cout++
		/*msg :=*/ client.Do(req)
		// fmt.Println(msg.Status, string(msg.ResBody))
		// if cout > 10000 {
		// 	break
		// }
		if time.Now().After(t) {
			fmt.Println("cout:", cout)
			return
		}
	}
	//==================================================
	// for i := 0; i < 10; i++ { //可跑1600*10个每秒
	// 	go func() {
	// 		cursor := client.Clone()
	// 		req := gcelery.NewTaskResquest()
	// 		req.Method = defaultMethod
	// 		req.ReqBody = []byte("你好")
	// 		req.Kwargs = map[string]string{
	// 			"aaa": "111",
	// 			"bbb": "222",
	// 		}
	// 		cout := 0
	// 		t := time.Now().Add(time.Second)
	// 		for { //单秒4200左右
	// 			cout++
	// 			/*msg :=*/ cursor.Do(req)
	// 			// fmt.Println(msg.Status, string(msg.ResBody))
	// 			// if cout > 10000 {
	// 			// 	break
	// 			// }
	// 			if time.Now().After(t) {
	// 				fmt.Println("cout:", cout)
	// 				return
	// 			}
	// 		}
	// 	}()
	// }
	time.Sleep(time.Second * 5)
	//====================================================

}

//go test
func main() {
	TestPool1()
	// TestClient1()

	//STL
	// TestPool2()
	// TestClient2()
	// TestTume()
}

//连接池+普通版本
func TestTume() {

	// addr := "localhost:50051"
	// defaultMethod := "Do1"
	// cursor := gcelery.NewCurSor(addr)

	// req := gcelery.NewTaskResquest()
	// req.Method = defaultMethod
	// req.ReqBody = []byte("你好")
	// req.Kwargs = map[string]string{
	// 	"aaa": "111",
	// 	"bbb": "222",
	// }
	// for i := 0; i < 100; i++ {
	// 	t := time.Now()
	// 	cursor.Do(req)
	// 	fmt.Println(time.Since(t))
	// }

	// fmt.Println(msg.Status, string(msg.ResBody))

}
