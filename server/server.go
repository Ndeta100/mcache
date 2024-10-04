package server

import (
	"fmt"
	"net"

	"github.com/Ndeta100/config"
)

type Server struct {
	cfg        config.Config
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
}

func NewServer() *Server {
	// Initialize the Server instance
	s := &Server{
		cfg:    config.GetDefaultConfig(), // Or use config.GetDefaultConfig()
		quitch: make(chan struct{}),
	}
	// Now that s is initialized, you can access s.cfg
	s.listenAddr = fmt.Sprintf("%s:%d", s.cfg.Network.Host, s.cfg.Network.Port)
	return s
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln
	go s.acceptLoop()
	<-s.quitch
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error", err)
			continue
		}
		fmt.Println("New connection to the server:", conn.RemoteAddr())
		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	buf := make([]byte, 2048)
	defer conn.Close()
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("reade error", err)
			continue
		}
		msg := buf[:n]
		fmt.Println(string(msg))
	}
}
