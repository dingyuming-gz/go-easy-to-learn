package main

import (
	"context"
	"fmt"
	pb "grpc-demo/server/proto"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//WithTransportCredentials返回一个DialOption，用于配置连接级别的安全凭证(例如，TLS/SSL)。这不能与WithCredentialsBundle一起使用。
	conn, err := grpc.Dial("127.0.0.1:9099", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect:%v", err)
	}
	defer conn.Close()
	//建立连接
	client := pb.NewSayHelloClient(conn)
	//rpc方法调用
	resp, _ := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: "client"})
	fmt.Println(resp.GetResponseMsg())
}
