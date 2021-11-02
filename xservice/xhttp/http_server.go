package xhttp

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/kucjac/cleango/pkg/xlog"
	"github.com/kucjac/cleango/xservice"
)

var _ xservice.RunnerCloser = (*Server)(nil)

// Server is a wrapper over http.Server that implements xservice.RunnerCloser interface.
type Server struct {
	Server *http.Server
}

// Run starts running the server in it's own go routine.
func (h *Server) Run() error {
	xlog.Printf("HTTP server listening at: %s", h.Server.Addr)
	go func() {
		err := h.Server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			xlog.Fatalf("Err: %v", err)
		}
	}()
	return nil
}

// Close stops the http server.
func (h *Server) Close(_ context.Context) error {
	return h.Server.Close()
}

// NewHTTPServer creates a new http server port for the images.
func NewHTTPServer(mux http.Handler, addr string) *Server {
	s := &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    time.Second * 30,
		WriteTimeout:   time.Second * 30,
		MaxHeaderBytes: 1 << 20,
	}
	return &Server{Server: s}
}
