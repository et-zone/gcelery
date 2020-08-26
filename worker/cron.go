package worker

import (
	"fmt"
	"log"
)

//定时任务，不需要传入参数，直接开启即可
func CronDo() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("cron err")
		}
	}()
	fmt.Println("cron test")
	panic("hhhh")
}
