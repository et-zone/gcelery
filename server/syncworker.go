package server

import (
	"golang.org/x/sync/errgroup"
)

var (
	g          errgroup.Group
	syncWorker *SyncWorker
)

type SyncWorker struct {
	funcs []func()
}

//异步持续运行的服务对象
func InitSyncWroker() *SyncWorker {
	if syncWorker != nil {
		return syncWorker
	}
	SyncWorker := &SyncWorker{}
	SyncWorker.funcs = []func(){}
	syncWorker = SyncWorker
	return syncWorker
}

func (this *SyncWorker) Init(fs ...func()) {

}

func (this *SyncWorker) RegisterWrokers(fs ...func()) {
	for _, f := range fs {
		this.funcs = append(this.funcs, f)
	}
}

func (this *SyncWorker) start() {
	for _, fun := range this.funcs {
		go fun()
	}

}

func StartSync(syncWorker *SyncWorker) {
	syncWorker.start()

}
