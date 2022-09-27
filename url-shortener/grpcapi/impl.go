package grpcapi

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
	v1 "url-shortener/grpcapi/pb/v1"
	"url-shortener/urlshortener"
)

func RunGRPCServer(manager urlshortener.Manager) {
	handler := &grpcHandler{manager: manager}

	server := grpc.NewServer()
	reflection.Register(server)
	v1.RegisterUrlShortenerServer(server, handler)

	lis, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalf("Failed to listen for gRPC connections due to error: %v", err)
	}

	log.Printf("Start serving gRPC at %q", lis.Addr())
	err = server.Serve(lis)
	log.Fatal(err)
}

type grpcHandler struct {
	v1.UnimplementedUrlShortenerServer

	manager urlshortener.Manager
}

func (h *grpcHandler) CreateShortcut(ctx context.Context, req *v1.CreateShortcutRequest) (*v1.CreateShortcutResponse, error) {
	fullURL := req.GetFullUrl()
	if fullURL == "" {
		return nil, status.Errorf(codes.InvalidArgument, "full URL must not be empty")
	}
	key, err := h.manager.CreateShortcut(ctx, fullURL)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create shortcut due to an error: %v", err)
	}

	return &v1.CreateShortcutResponse{
		Key: key,
	}, nil
}
