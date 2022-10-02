package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"speed-up/speedup"
	"time"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := speedup.NewDataServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	response, err := client.SetData(ctx, &speedup.RequestDataKeyValue{
		Key:   "nome",
		Value: "Thiago joses",
	})

	if err != nil {
		log.Fatalf("%v", err)
	}

	if response.GetException() != "" {
		log.Fatalf(response.GetException())
	}

	resp, err := client.GetData(ctx, &speedup.RequestDataKey{
		Key: "nome",
	})

	if err != nil {
		log.Fatalf("%v", err)
	}

	if response.GetException() != "" {
		log.Fatalf(response.GetException())
	}

	fmt.Printf("Resposta %v", resp.GetValue())
}
