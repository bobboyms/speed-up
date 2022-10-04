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
	conn, err := grpc.Dial("127.0.01:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := speedup.NewDataServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), time.Second)

	response, _ := client.SetData(ctx, &speedup.RequestDataKeyValue{
		Key:   "nome",
		Value: "Thiago joses",
	})

	response, _ = client.SetData(ctx, &speedup.RequestDataKeyValue{
		Key:   "idade",
		Value: "13",
	})

	if response.GetException() != "" {
		log.Fatalf(response.GetException())
	}

	requests := []*speedup.RequestDataKey{{Key: "idade"}, {Key: "nome"}}

	data := &speedup.RequestDataKeyList{
		RequestDataKeyList: requests,
	}

	resps, _ := client.GetsData(ctx, data)

	println(resps.ResponseDataValueList[1].Value)

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
