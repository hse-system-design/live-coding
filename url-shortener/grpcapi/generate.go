package grpcapi

// requires to run the following:
// brew install protobuf
// go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
// go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
// also make sure that your $GOPATH is listed in your $PATH
//go:generate protoc -I . --go_out=. --go-grpc_out=. urlshortener.proto
