package server

import (
	"errors"
	// "fmt"
	"reflect"
	"runtime"
	"strings"
)

var globalWorker *RpcWorker

type RpcWorker struct {
	workerMap map[string]func([]byte) (error, []byte)
}

func NewRpcWorker() {
	if globalWorker != nil {
		return
	}
	worker := &RpcWorker{}
	worker.workerMap = make(map[string]func([]byte) (error, []byte), 0)
	globalWorker = worker
}

func RegisterRpcWorker(fs ...func([]byte) (error, []byte)) error {
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

func GetRpcWorker(name string) func([]byte) (error, []byte) {
	return globalWorker.workerMap[name]
}
