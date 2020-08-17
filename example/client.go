package main

import (
	"fmt"
	"strings"
	"time"

	// "testing"
	// "context"
	// "time"
	"flag"
	"strconv"

	"github.com/et-zone/gcelery"
)

//连接池+普通版本
func TestPool1() { // -t 10
	flg := flag.String("t", "20", "please input goroutines count!")
	flag.Parse()
	size, _ := strconv.Atoi(*flg)
	fmt.Println("worker :", size)
	addr := "172.17.0.14:50051"
	// addr := "172.17.0.3:50051"
	// addr := "49.232.190.114:50051"
	defaultMethod := "Do1"
	client := gcelery.NewClient(addr)
	defer client.Close()

	req := gcelery.NewTaskResquest()
	req.Method = defaultMethod
	req.ReqBody = []byte("你好")
	req.Kwargs = map[string]string{
		"aaa": "111",
		"bbb": "222",
	}
	for i := 0; i < size; i++ {
		go func() {
			count := 0

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
	client := gcelery.NewSTLClient(addr, "/tl/server.crt")
	defer client.Close()
	req := gcelery.NewTaskResquest()
	req.Method = defaultMethod
	req.ReqBody = []byte("你好")
	req.Kwargs = map[string]string{
		"aaa": "111",
		"bbb": "222",
	}
	cout := 0
	t := time.Now().Add(time.Second * 1)
	for {

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
	client := gcelery.NewClient(addr)
	defer client.Close()

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
	client := gcelery.NewSTLClient(addr, "/tl/server.crt")
	defer client.Close()

	req := gcelery.NewTaskResquest()
	req.Method = defaultMethod
	req.ReqBody = []byte("你好")
	req.Kwargs = map[string]string{
		"aaa": "111",
		"bbb": "222",
	}
	cout := 0
	t := time.Now().Add(time.Second)
	for {
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
	time.Sleep(time.Second * 5)

}

//go test
func main() {
	TestPool1()
	// TestClient1()

	//STL
	// TestPool2()
	// TestClient2()
}
