package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"id-maker/internal/controller/rpc/proto"
	"log"
)

const (
	// Address 连接地址
	Address string = ":50051"
)

func main() {
	// 连接服务器
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("net.Connect err: %v", err)
	}
	defer conn.Close()

	// 建立gRPC连接
	grpcClient := proto.NewGidClient(conn)
	// 创建发送结构体
	req := proto.IdRequest{
		Tag: "test",
	}
	//req := emptypb.Empty{}
	// 调用我们的服务(Route方法)
	// 同时传入了一个 context.Context ，在有需要时可以让我们改变RPC的行为，比如超时/取消一个正在运行的RPC

	res, err := grpcClient.GetId(context.Background(), &req)
	if err != nil {
		log.Fatalf("Call Route err: %v", err)
	}
	// 打印返回值
	log.Println(res.Id)
}
