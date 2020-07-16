package server

import (
	cron "gopkg.in/robfig/cron.v2"
)

type Cron struct {
	cron *cron.Cron
}

var cro *Cron

func NewCron() *Cron {
	if cro != nil {
		return cro
	}
	cro = &Cron{cron: cron.New()}
	return cro
}

func (this *Cron) RegisterWroker(spec string, cmd func()) (cron.EntryID, error) {
	return this.cron.AddFunc(spec, cmd)
}

func (this *Cron) start() {
	this.cron.Start()

}

func StartCron(cron *Cron) {
	cron.start()

}
func (this *Cron) stop() {
	this.cron.Stop()
}
func StopCron(cron *Cron) {
	cron.stop()
}

//用不上
// func (this *Cron) Select() {
// 	select {}
// }
