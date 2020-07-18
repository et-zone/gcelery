package control

import (
	"errors"
	// "fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/et-zone/gcelery/task"
)

var globalWorker *CeleryWorker

type CeleryWorker struct {
	workerMap map[string]func(*task.Request) (error, *task.Response)
}

func NewCeleryWorker() {
	if globalWorker != nil {
		return
	}
	worker := &CeleryWorker{}
	worker.workerMap = make(map[string]func(*task.Request) (error, *task.Response), 0)
	globalWorker = worker
}

func RegisterCeleryWorker(fs ...func(*task.Request) (error, *task.Response)) error {
	if fs == nil {
		return errors.New("RegisterWorker err,variable is null ")
	}
	for _, f := range fs {
		p := reflect.ValueOf(f).Pointer()
		list := strings.SplitAfter(runtime.FuncForPC(p).Name(), ".")
		globalWorker.workerMap[list[len(list)-1]] = f
	}
	// fmt.Println(globalWorker.workerMap)
	return nil

}

func DeRegisterRpc(f func([]byte) (error, []byte)) {
	if f == nil {
		return
	}
	p := reflect.ValueOf(f).Pointer()
	list := strings.SplitAfter(runtime.FuncForPC(p).Name(), ".")
	delete(globalWorker.workerMap, list[len(list)-1])
}

func GetCeleryWorker(name string) func(*task.Request) (error, *task.Response) {
	return globalWorker.workerMap[name]
}
