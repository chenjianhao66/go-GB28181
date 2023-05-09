package syn

import (
	"github.com/chenjianhao66/go-GB28181/internal/pkg/log"
	"github.com/pkg/errors"
	"time"
)

// todo 等待重构
type Entity struct {
	key    string
	err    chan error
	data   chan interface{}
	ticker *time.Timer
}

var (
	ErrTimeOut = errors.New("time out")
)

func (e *Entity) Wait() (interface{}, error) {
	defer e.destroy()
	select {
	case err := <-e.err:
		log.Error(err)
		return nil, err
	case d := <-e.data:
		return d, nil
	case <-e.ticker.C:
		return nil, ErrTimeOut
	}
}

func (e *Entity) destroy() {
	d.mux.Lock()
	defer d.mux.Unlock()
	delete(d.d, e.key)
	_ = e.Close()
}

func (e *Entity) Close() error {
	close(e.data)
	e.ticker.Stop()
	close(e.err)
	return nil
}

func (e *Entity) Ok(data interface{}) {
	e.data <- data
}

func (e *Entity) Err(err error) {
	e.err <- err
}
