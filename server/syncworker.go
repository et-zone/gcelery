package server

import (
	"golang.org/x/sync/errgroup"
)

var (
	g          errgroup.Group
	syncWroker *SyncWroker
)

type SyncWroker struct {
	funcs []func()
}

//异步持续运行的服务对象
func InitSyncWroker() *SyncWroker {
	if syncWroker != nil {
		return syncWroker
	}
	SyncWroker := &SyncWroker{}
	SyncWroker.funcs = []func(){}
	syncWroker = SyncWroker
	return syncWroker
}

func (this *SyncWroker) Init(fs ...func()) {

}

func (this *SyncWroker) RegisterWrokers(fs ...func()) {
	for _, f := range fs {
		this.funcs = append(this.funcs, f)
	}
}

func (this *SyncWroker) start() {
	for _, fun := range this.funcs {
		go fun()
	}

}

func StartSync(syncWroker *SyncWroker) {
	syncWroker.start()

}
