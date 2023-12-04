package client

import (
	"bytes"
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework/internal/grpc/generated"
	"io"
	"time"
)

type FileService struct {
	client generated.FileServiceClient
	cfg    Config
	log    *zap.Logger
}

type Config struct {
	timeout time.Duration
	uri     string
}

func New(cfg Config, log *zap.Logger) (*FileService, error) {
	conn, err := grpc.Dial(cfg.uri, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	client := generated.NewFileServiceClient(conn)
	return &FileService{
		client: client,
		cfg:    cfg,
		log:    log,
	}, nil
}

func (s *FileService) GetFile(ctx context.Context, req *generated.FileRequest) (*File, error) {
	ctx, cancel := context.WithTimeout(ctx, s.cfg.timeout)
	defer cancel()

	c, err := s.client.GetFile(ctx, req)
	if err != nil {
		return nil, err
	}

	res := File{}
	ch := make(chan []byte)
	errCh := make(chan error)
	buf := bytes.Buffer{}
	go func() {
		for {
			resp, err := c.Recv()
			if err == io.EOF {
				close(ch)
				break
			}
			if err != nil {
				errCh <- err
			}

			res.Name = resp.Name
			res.Size = resp.Size
			ch <- resp.Content
		}
	}()

	select {
	case content, ok := <-ch:
		if !ok {
			break
		}
		buf.Write(content)
	case <-ctx.Done():
		return nil, ctx.Err()
	case err = <-errCh:
		return nil, err
	}

	fullContent := make([]byte, buf.Len())
	_, err = buf.Read(fullContent)
	if err != nil {
		return nil, err
	}

	res.Content = fullContent
	return &res, nil
}

type File struct {
	Name    string
	Size    int32
	Content []byte
}
