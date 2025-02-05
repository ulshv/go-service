package server

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
)

type Server struct {
	address string
	port    string
	server  *http.Server
}

type NewServerOptions struct {
	Address string
	Port    string
}

func NewServer(opts NewServerOptions) *Server {
	mux := http.NewServeMux()
	// routes
	mux.HandleFunc("/hello", HelloHandler)
	mux.HandleFunc("/foo", FooHandler)
	mux.HandleFunc("/bar", BarHandler)
	// server instance
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", opts.Address, opts.Port),
		Handler: mux,
	}
	return &Server{
		address: opts.Address,
		port:    opts.Port,
		server:  server,
	}
}

func (s *Server) Listen() {
	slog.Info("Server is starting listenting", "address", s.address, "port", s.port)

	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
