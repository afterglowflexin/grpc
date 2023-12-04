package server

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"homework/internal/grpc/generated"
	"net"
)

type Server struct {
	generated.UnimplementedFileServiceServer
	app application
	cfg Config
	log *zap.Logger
}

type Config struct {
	address string
}

type application interface {
	GetFile(name string)
	GetFilesNames()
	GetFileInfo()
}

func New(cfg Config, app application, log *zap.Logger) *Server {
	return &Server{
		UnimplementedFileServiceServer: generated.UnimplementedFileServiceServer{},
		app:                            app,
		cfg:                            cfg,
		log:                            log,
	}
}

func (s *Server) ListenAndServe() error {
	lis, err := net.Listen("tcp", s.cfg.address)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	generated.RegisterFileServiceServer(server, &Server{})
	if err := server.Serve(lis); err != nil {
		return err
	}

	return nil
}
