package main

import (
	// "testing"

	do "github.com/et-zone/gcelery/wroker" //yourself wroker

	"github.com/et-zone/gcelery"
)

//普通版本
func TestServer() {
	addr := ":50051"
	//初始化服务器
	server := gcelery.NewCelery(addr) //创建服务器对象
	//初始化worker //业务==grpc
	server.InitCelery()
	//注册worker
	server.RegisterCeleryWorker(do.Do1)
	//注册通信协议服务
	server.RegisterTransport()
	//业务2==cron
	//针对定时任务，定时任务不需要返回消息
	cron := server.NewCronWorker()
	cron.RegisterWroker("*/4 * * * * *", do.CronDo)
	//异步运行的服务
	swroker := server.NewSyncWroker()
	swroker.RegisterWrokers(do.LongDo, do.LongDo2)
	//启动运行服务
	// server.RegisterCron(cron)
	// server.RegisterSync(swroker)
	server.StartCelery() //启动服务

}

//STL 版本
func TestServerSTL() {
	addr := ":50051"
	//初始化服务器
	server := gcelery.NewTlsCelery(addr, "/tl/server.crt", "/tl/server.key") //创建服务器对象
	//初始化worker //业务==grpc
	server.InitCelery()
	//注册worker
	server.RegisterCeleryWorker(do.Do1)
	//注册通信协议服务
	server.RegisterTransport()
	//业务2==cron
	//针对定时任务，定时任务不需要返回消息
	cron := server.NewCronWorker()
	cron.RegisterWroker("*/4 * * * * *", do.CronDo)
	//异步运行的服务
	swroker := server.NewSyncWroker()
	swroker.RegisterWrokers(do.LongDo, do.LongDo2)
	//启动运行服务
	// server.RegisterCron(cron)
	// server.RegisterSync(swroker)
	server.StartCelery() //启动服务

}

func main() {
	TestServer()
	// TestServerSTL()
}
