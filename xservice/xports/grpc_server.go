package xports

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/kucjac/cleango/cgerrors"
	"github.com/kucjac/cleango/xlog"
	"github.com/kucjac/cleango/xservice"
)

// Compile time check if the GRPCServer implements xservice.RunnerCloser interface.
var _ xservice.RunnerCloser = (*GRPCServer)(nil)

// GRPCServer is a wrapper over grpc.Server that implements xservice.RunnerCloser interface.
type GRPCServer struct {
	*grpc.Server
	Listener net.Listener
}

// Run starts the GRPC server.
func (g *GRPCServer) Run() error {
	go func() {
		err := g.Server.Serve(g.Listener)
		if err != nil && !cgerrors.Is(err, grpc.ErrServerStopped) {
			xlog.Fatalf("Err: %v", err)
		}
	}()
	xlog.Infof("Listening GRPC on: %s", g.Listener.Addr())
	return nil
}

// Close stops the GRPC server.
func (g *GRPCServer) Close(_ context.Context) error {
	g.Server.GracefulStop()
	return nil
}

// NewGRPCServer creates a new port for the GRPC server.
func NewGRPCServer(addr string, options ...grpc.ServerOption) (*GRPCServer, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	// Create new GRPC server with provided options.
	s := grpc.NewServer(options...)
	srv := &GRPCServer{Server: s, Listener: lis}
	return srv, nil
}
