package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/Ndeta100/mcache/config"
	"github.com/Ndeta100/mcache/handler"
	"github.com/Ndeta100/mcache/store"
)

type Server struct {
	cfg        config.Config
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	cache      *store.Cache
}

func NewServer(cache *store.Cache) *Server {
	// Initialize the Server instance
	s := &Server{
		cfg:    config.GetDefaultConfig(), // Or use config.GetDefaultConfig()
		quitch: make(chan struct{}),
		cache:  cache,
	}
	// Now that s is initialized, you can access s.cfg
	s.listenAddr = fmt.Sprintf("%s:%d", s.cfg.Network.Host, s.cfg.Network.Port)
	return s
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer ln.Close()
	s.ln = ln

	fmt.Printf("Server started on %s\n", s.listenAddr)
	go s.acceptLoop()
	<-s.quitch // Block until quit signal is received
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}
		fmt.Println("New connection to the server:", conn.RemoteAddr())
		go s.readLoop(conn) // Start a Goroutine to handle each connection
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		// Read until a newline character is found
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Read error:", err)
			return // Exit the loop if there's an error reading from the connection
		}

		// Trim spaces and newline characters
		command = strings.TrimSpace(command)

		// Handle the command using the handler
		response := handler.HandleCommand(command, s.cache)

		// Write the response back to the client with proper newline for telnet
		_, writeErr := conn.Write([]byte(response + "\r\n"))
		if writeErr != nil {
			fmt.Printf("Write error: %v\n", writeErr)
			return // Exit if unable to write response
		}
	}
}
