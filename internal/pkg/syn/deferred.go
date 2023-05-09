package syn

import (
	"sync"
	"time"
)

type deferredSync struct {
	mux *sync.RWMutex
	d   map[string]*Entity
}

var d = &deferredSync{
	mux: &sync.RWMutex{},
	d:   make(map[string]*Entity),
}

func NewDelayTask(key string, timeout time.Duration) *Entity {
	d.mux.Lock()
	defer d.mux.Unlock()
	e := &Entity{
		key:    key,
		err:    make(chan error, 1),
		data:   make(chan interface{}, 1),
		ticker: time.NewTimer(timeout),
	}
	d.d[key] = e
	return e
}

type CallBackFunc func(e *Entity)

func HasSyncTask(key string, successCall ...CallBackFunc) bool {
	d.mux.RLock()
	defer d.mux.RUnlock()
	e, ok := d.d[key]
	if len(successCall) > 0 && ok {
		for _, callBackFunc := range successCall {
			callBackFunc(e)
		}
	}
	return ok
}
