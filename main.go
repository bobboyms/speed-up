package main

import (
	"context"
	"fmt"
	"github.com/dustin/go-humanize"
	"google.golang.org/grpc"
	"log"
	"net"
	"runtime"
	"speed-up/service/application/ports"
	"speed-up/service/application/service"
	"speed-up/service/infraestructure/memory"
	"speed-up/service/infraestructure/zaplog"
	"speed-up/speedup"
	"time"
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

func (s *Server) GetsData(ctx context.Context, req *speedup.RequestDataKeyList) (*speedup.ResponseDataValueList, error) {

	keys := make([]string, len(req.RequestDataKeyList))
	for i, requestData := range req.RequestDataKeyList {
		keys[i] = requestData.Key
	}

	values, err := s.Service.Gets(keys...)

	if err != nil {
		return &speedup.ResponseDataValueList{
			Exception: err.Error(),
		}, nil
	}

	responses := make([]*speedup.ResponseDataValue, len(values))
	for i, value := range values {
		responses[i] = &speedup.ResponseDataValue{
			Value: value,
		}
	}

	return &speedup.ResponseDataValueList{
		ResponseDataValueList: responses,
	}, nil

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

func info() {

	for {
		var startMemory runtime.MemStats
		runtime.ReadMemStats(&startMemory)

		fmt.Printf("============= STATUS MEMORY ===============\n")
		fmt.Printf("Memory loc %v\n", humanize.Bytes(startMemory.Alloc))
		fmt.Printf("Number CPUs %v\n", runtime.NumCPU())
		fmt.Printf("%d goroutines running\n", runtime.NumGoroutine())
		time.Sleep(5 * time.Second)
	}

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
	go info()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
