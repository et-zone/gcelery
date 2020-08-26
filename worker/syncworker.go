package worker

import (
	"fmt"
	"log"
	"time"
)

var count = 0

//异步永久执行的程序，不需要参数，如果需要可走mq或者db查询
func LongDo() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("LongDo err 要重启服务了 不正常退出时会重启", err)
			go LongDo()
		}
		fmt.Println("程序正常退出")

	}()
	for {
		count++
		time.Sleep(time.Second)
		fmt.Println("long do! ", count)

		if count > 7 { //正常退出==此时程序会退出
			return
		}
		if count/5 > 0 {
			panic("sadf") //会走go程序，继续纸箱
		}

	}

}

func LongDo2() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("LongDo2 err")
		}
	}()

	for {
		time.Sleep(time.Second)
		fmt.Println("long do2!")

	}

}
