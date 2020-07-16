package server

import (
	"errors"
	// "fmt"
	"context"
	"reflect"
	"runtime"
	"strings"
)

var globalWorker *CeleryWorker

type CeleryWorker struct {
	workerMap map[string]func(context.Context, []byte) (error, []byte)
}

func NewCeleryWorker() {
	if globalWorker != nil {
		return
	}
	worker := &CeleryWorker{}
	worker.workerMap = make(map[string]func(context.Context, []byte) (error, []byte), 0)
	globalWorker = worker
}

func RegisterRpcWorker(fs ...func(context.Context, []byte) (error, []byte)) error {
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

func GetRpcWorker(name string) func(context.Context, []byte) (error, []byte) {
	return globalWorker.workerMap[name]
}
