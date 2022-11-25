package sip

import (
	"github.com/chenjianhao66/go-GB28181/log"
	"net"
)

type Server struct {
	tcp handler
	udp handler
	p   *parser
}

func NewServer() *Server {
	p := newParser()
	return &Server{
		p: p,
	}
}

func (s *Server) Run() {
	go s.listenUdpServer()
	//go s.listenTcpServer()
	select {}
}

func (s *Server) listenUdpServer() {
	log.Log.Info("enter listenUdpServer function....")
	// TODO 将地址放到配置文件中
	addr, err := net.ResolveUDPAddr("udp", "192.168.1.223:5060")
	if err != nil {
		log.Log.Fatalf("resolve fail from %s: \n %s", "192.168.1.223:5060", err)
	}

	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Log.Fatalf("listen upd message fail from %s : \n %s", "192.168.1.223:5060", err)
	}
	buf := make([]byte, BufferSizeMax)
	for {
		num, udpAddr, err := udpConn.ReadFrom(buf)
		if err != nil {
			log.Log.Fatalf("read udp message from udpConn : \n%s", err)
		}
		log.Log.Debugf("receive from remote addr %s and %d size upd message from updConn", udpAddr, num)
		log.Log.Debugf("receive the data : \n%s", string(buf[:num]))
		s.p.in <- newPacket(append([]byte{}, buf[:num]...), udpAddr)
		buf = buf[:]
	}

}

func (s *Server) listenTcpServer() {
	// TODO 将地址放到配置文件中
	//listener, err := net.Listen("tcp", "192.168.1.223:5060")
	//if err != nil {
	//	log.Log.Error("resolve tcp addr fail")
	//	panic(any(err))
	//}
	//
	//// listen tcp protocol sip message
	//for {
	//	conn, err := listener.Accept()
	//	if err != nil {
	//		log.Log.Error("error conn,break")
	//		break
	//	}
	//	go func() {
	//		s.tcp.Handle(conn,s.p.in)
	//	}()
	//
	//}

}

func (s *Server) Close() error {
	_ = s.tcp.Close()
	_ = s.udp.Close()
	s.p.stop()
	return nil
}
