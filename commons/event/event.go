package event

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	HightPriority = iota
	MiddlePriority
	LowPriority
	Max
)

var (
	hookFuncs [][]func()
	hook      chan os.Signal
	mux       *sync.Mutex
)

func init() {
	hookFuncs = make([][]func(), Max)
	hook = make(chan os.Signal, 1)
	mux = new(sync.Mutex)
}

func AddHook(priority int, f func()) bool {
	if priority < 0 || priority >= Max {
		return false
	}

	mux.Lock()
	funcs := hookFuncs[priority]
	funcs = append(funcs, f)
	hookFuncs[priority] = funcs
	mux.Unlock()
	return true
}

// 优先级从高到底执行
func doHook() {
	for priority := 0; priority < Max; priority++ {
		funcs := hookFuncs[priority]
		for _, hookFunc := range funcs {
			hookFunc()
		}
	}
}

// 等待信号
func Wait(signals ...os.Signal) {
	if len(signals) == 0 {
		signal.Notify(hook, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGSTOP)
	} else {
		signal.Notify(hook, signals...)
	}
	<-hook
}

func Exit(code int) {
	os.Exit(code)
}

func WaitExit(signals ...os.Signal) {
	Wait(signals...)
	doHook()
	Exit(0)
}

// 暂时不提供停止方法
// func Stop() {
// 	signal.Stop(hook)
// 	close(hook)
// }
