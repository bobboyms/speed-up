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
	//ec2-34-239-251-75.compute-1.amazonaws.com:9000
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("127.0.0.1:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := speedup.NewDataServiceClient(conn)
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

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

	requests := []*speedup.RequestDataKey{{Key: "idadex"}, {Key: "nomex"}, {Key: "xxx"}, {Key: "bbb"}}

	data := &speedup.RequestDataKeyList{
		RequestDataKeyList: requests,
	}

	resps, _ := client.GetsData(ctx, data)

	for _, value := range resps.ResponseDataValueList {
		println(value.Value)
	}

	resp, err := client.GetData(ctx, &speedup.RequestDataKey{
		Key: "idade",
	})

	if err != nil {
		log.Fatalf("%v", err)
	}

	println("XXX")
	if response.GetException() != "" {
		log.Fatalf(response.GetException())
	}

	fmt.Printf("Resposta %v", resp.Found)
}
