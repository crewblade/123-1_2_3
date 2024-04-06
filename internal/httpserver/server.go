package httpserver

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(
	router http.Handler,
	readTimeout time.Duration,
	writeTimeout time.Duration,
	shutdownTimeout time.Duration,
	addr string,
) *Server {
	httpServer := &http.Server{
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Addr:         addr,
	}

	server := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: shutdownTimeout,
	}

	server.start()

	return server
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
