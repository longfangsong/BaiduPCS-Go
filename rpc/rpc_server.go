package rpc

import (
	"context"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type Server struct {
	UnimplementedFileServiceServer
	client cli.App
}

func NewServer(client cli.App) Server {
	return Server{
		UnimplementedFileServiceServer: UnimplementedFileServiceServer{},
		client:                         client,
	}
}
func (s *Server) Download(ctx context.Context, req *DownloadRequest) (*DownloadResponse, error) {
	s.client.Run([]string{"BaiduPCS-Go", "download", "/lana/" + req.Filename, "--saveto", "/download", "--nocheck", "-l", "3"})
	file, _ := os.Open("/download/" + req.Filename)
	defer file.Close()
	content, _ := ioutil.ReadAll(file)
	return &DownloadResponse{
		Data: content,
	}, nil
}
func (s *Server) Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	file, _ := os.Create("/upload/" + req.Filename)
	defer file.Close()
	file.Write(req.Data)
	s.client.Run([]string{"BaiduPCS-Go", "upload", "/upload/" + req.Filename, "/lana"})
	return &UploadResponse{
		Success: true,
	}, nil
}
func (s *Server) Listen(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	RegisterFileServiceServer(server, s)
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
