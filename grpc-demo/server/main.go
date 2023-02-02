package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "grpc-demo/server/proto"

	"google.golang.org/grpc"
)

// hello server 服务器端
type server struct {
	pb.UnimplementedSayHelloServer
}

// 方法重写
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{ResponseMsg: "hello " + req.RequestName}, nil
}

func main() {
	//开启端口
	listen, err := net.Listen("tcp", ":9099")
	if err != nil {
		fmt.Println(err)
	}
	//创建grpc服务
	grpcServer := grpc.NewServer()
	//在grpc服务端中去注册我们自己编写的服务,注册一定是引用注册
	pb.RegisterSayHelloServer(grpcServer, &server{})
	//启动服务
	err = grpcServer.Serve(listen)
	if err != nil {
		log.Println("启动失败~")
	}
}
