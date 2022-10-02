package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"speed-up/service/application/ports"
	"speed-up/service/application/service"
	"speed-up/service/infraestructure/memory"
	"speed-up/service/infraestructure/zaplog"
	"speed-up/speedup"
)

type Server struct {
	speedup.UnimplementedDataServiceServer
	Service ports.MemoryData
}

func NewServer(service ports.MemoryData) speedup.DataServiceServer {
	return &Server{
		Service: service,
	}
}

func (s *Server) GetData(ctx context.Context, req *speedup.RequestDataKey) (*speedup.ResponseDataValue, error) {

	value, err := s.Service.Get(req.GetKey())

	if err != nil {
		return &speedup.ResponseDataValue{
			Value:     "",
			Exception: err.Error(),
		}, err
	}

	return &speedup.ResponseDataValue{
		Value:     value,
		Exception: "",
	}, nil

}
func (s *Server) SetData(ctx context.Context, req *speedup.RequestDataKeyValue) (*speedup.ResponseEmpty, error) {
	err := s.Service.Set(req.GetKey(), req.GetValue())
	if err != nil {
		return &speedup.ResponseEmpty{
			Exception: err.Error(),
		}, err
	}

	return &speedup.ResponseEmpty{}, nil
}

func main() {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logger := zaplog.NewZapLogger()
	repository := memory.NewMemoryIndexRepository(logger)
	service := service.NewMemoryService(repository)

	grpcServer := grpc.NewServer()
	server := NewServer(service)

	speedup.RegisterDataServiceServer(grpcServer, server)

	log.Printf("GRPC server listening on %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
