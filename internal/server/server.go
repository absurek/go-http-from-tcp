package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/absurek/go-http-from-tcp/internal/request"
	"github.com/absurek/go-http-from-tcp/internal/response"
)

type HandlerError struct {
	Status  response.StatusCode
	Message string
}

func (e *HandlerError) Dump(w io.Writer) {
	err := response.WriteStatusLine(w, e.Status)
	if err != nil {
		log.Printf("Error: writing error status line: %v", err)
	}

	err = response.WriteHeaders(w, response.GetDefaultHeaders(len(e.Message)))
	if err != nil {
		log.Printf("Error: writing error headers: %v", err)
	}

	w.Write([]byte(e.Message))
}

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	listener net.Listener
	handler  Handler
}

func Serve(port int, handler Handler) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("setting up listener: %v", err)
	}

	server := &Server{listener: listener, handler: handler}
	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Printf("Error: accepting connection: %v", err)
			continue
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Printf("Error: creating request from connection: %v", err)
		return
	}

	var writer response.Writer
	s.handler(&writer, req)

	conn.Write(writer.Bytes())
}
