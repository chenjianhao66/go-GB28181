package sip

import (
	"io"
	"net"
)

type connection interface {
	io.ReadWriteCloser
	Network() string
}

type tcpConnection struct {
	baseConn   net.Conn
	localAddr  net.Addr
	remoteAddr net.Addr
}

type udpConnection struct {
	baseConn   net.Conn
	localAddr  net.Addr
	remoteAddr net.Addr
}
