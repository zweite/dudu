package mymgo

import (
	"errors"
	"sync"

	"dudu/config"
)

var (
	mux       sync.RWMutex
	once      sync.Once
	cacheSess map[string]*MdbSession
	resConfig *config.ResourceConfig
)

func init() {
	cacheSess = make(map[string]*MdbSession)
}

// if can't dial mongodb that it will be panic
func OpenCfg(mgoCfgs []*config.DBConfig) (sess *MdbSession) {
	sess = GetDb(mgoCfgs)
	return
}

func Open(name string) (sess *MdbSession, err error) {
	var ok bool
	mux.RLock()
	sess, ok = cacheSess[name]
	mux.RUnlock()
	if ok {
		return
	}

	once.Do(func() {
		resConfig = config.ParseResourceConfig()
		if resConfig == nil {
			panic("res config can't be nil")
		}
	})

	mgoCfgs, ok := resConfig.Mongo[name]
	if !ok {
		err = errors.New("配置不存在")
		return
	}

	sess = GetDb(mgoCfgs)
	mux.Lock()
	cacheSess[name] = sess
	mux.Unlock()
	return
}

func OpenBase() (db *MdbSession, err error) {
	return Open("base")
}
