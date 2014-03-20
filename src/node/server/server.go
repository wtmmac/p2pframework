package server

import (
	"net"
	"fmt"
)

type Server struct {
	listener net.Listener
	port int
	upstream chan interface{}
	running bool
}

type OnConnection struct {
	Connection net.Conn
}

func NewServer(upstream chan interface{}, port int) *Server {
	server := new(Server)

	server.upstream = upstream
	server.port = port

	return server
}

func (server *Server) Start() {
	go server.mainLoop()
}

func (server *Server) Stop() {
	server.listener.Close()
}

func (server *Server) mainLoop() {
	server.running = true

	var err error
	server.listener, err = net.Listen("tcp", ":8000")

	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		server.running = false
		return
	}

	for {
		conn, err := server.listener.Accept()

		if err == nil {
			server.upstream <- OnConnection{conn}
		} else {
			break
		}
	}

	server.running = false
}
