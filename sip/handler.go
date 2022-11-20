package sip

import (
	"net"
	"sync"
	"sync/atomic"
)

type handler interface {
	Handle(conn net.Conn)
	Close() error
}

type tcpHandler struct {
	// online conn
	activeConn sync.Map
	// when this handler is close status, that field is true
	closing atomic.Bool

	wg sync.WaitGroup
}

func (h *tcpHandler) Handle(conn net.Conn) {

}

func (h *tcpHandler) Close() error {
	h.activeConn.Range(func(key, value any) bool {
		conn := key.(net.Conn)
		_ = conn.Close()
		return true
	})
	h.closing.Store(true)
	return nil
}
