package main

import (
	"go.uber.org/zap"
	"homework/internal/application"
	"homework/internal/config"
	"homework/internal/grpc/server"
	"os"
)

type File struct {
	Name    string `json:"name"`
	Size    int32  `json:"size"`
	Content []byte `json:"content"`
}

func main() {
	cfg := config.Config{}
	app := application.New()
	log, _ := zap.NewProduction()

	server := server.New(cfg.Server, app, log)

	err := server.ListenAndServe()
	if err != nil {
		os.Exit(1)
	}
}
