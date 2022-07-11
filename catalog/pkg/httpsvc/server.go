package httpsvc

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Server represents an HTTP server
type Server struct {
	srv           *http.Server
	shutdownGrace time.Duration
}

// NewServer returns a new instance of server
func NewServer(handler http.Handler, opts ...ServerOption) (*Server, error) {
	srv := &http.Server{
		Addr:              ":3000", // TODO: Look into this
		Handler:           handler, // TODO: Look into this
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		// BaseContext:       nil, // TODO: Look into this
		// ConnContext:       nil, // TODO: Look into this
	}

	s := &Server{srv: srv, shutdownGrace: 10 * time.Second}

	for _, opt := range opts {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	startErrChan := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Println("Starting server...")
		startErrChan <- s.srv.ListenAndServe()
	}()

	// Blocking main and waiting for shutdown.
	select {
	case err := <-startErrChan:
		if err != http.ErrServerClosed { // ListenAndServe will always return a non-nil error
			return fmt.Errorf("startup failed: %w", err)
		}
		return nil
	case <-ctx.Done():
		return s.stop()
	}
}

func (s *Server) stop() error {
	log.Println("Shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownGrace)
	defer cancel()

	log.Println("Attempting graceful shutdown...")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Println(fmt.Sprintf("graceful shutdown did not complete in %v : %+v", s.shutdownGrace, err))

		log.Println("Attempting force shutdown...")
		if err = s.srv.Close(); err != nil {
			return fmt.Errorf("force shutdown failed: %w", err)
		}
	}
	log.Println("Server shutdown successfully!")

	return nil
}

// ServerOption is an optional config used to modify the server's behaviour
type ServerOption func(*Server) error

// WithServerPort overrides the server's default port with the given port
func WithServerPort(port string) ServerOption {
	return func(s *Server) error {
		if _, err := strconv.Atoi(port); err != nil {
			return fmt.Errorf("invalid port. err: %w", err)
		}

		s.srv.Addr = fmt.Sprintf(":%s", port)
		return nil
	}
}
