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
	cr := &Cron{cron: cron.New()}
	cro = cr
	return cro
}

func (this *Cron) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	return this.cron.AddFunc(spec, cmd)
}

func (this *Cron) Start() {
	this.cron.Start()

}
func (this *Cron) Stop() {
	this.cron.Stop()
}

//用不上
// func (this *Cron) Select() {
// 	select {}
// }
