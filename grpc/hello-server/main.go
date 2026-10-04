package main

import (
	"context"
	"fmt"
	pb "grpc/hello-server/proto"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	//获取元数据信息
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not found")
	}
	var appId string
	var appKey string
	if v, ok := md["appid"]; ok {
		appId = v[0]
	}
	if v, ok := md["appkey"]; ok {
		appKey = v[0]
	}
	if appId != "wjh" || appKey != "1234567890" {
		return nil, status.Errorf(codes.Unauthenticated, "appId or appKey is not valid")
	}

	return &pb.HelloResponse{
		Returnmessage: "Hello, " + req.GetName(),
	}, nil
}

func main() {
	//开启端口
	listen, err := net.Listen("tcp", ":9090")
	if err != nil {
		fmt.Println("开启端口失败,err:", err)
		return
	}
	defer listen.Close()

	//创建grpc服务器
	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	//在grpc服务器中注册服务
	pb.RegisterHelloServiceServer(grpcServer, &server{})

	//启动服务器
	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Println("启动服务器失败,err:", err)
		return
	}
}

//服务端编译命令：protoc --go_out=. --go-grpc_out=. proto/hello.proto
//服务端编写流程：
//1，创建gRPC Sever对象
//2，在将server(包含需要被调用的服务器接口)注册到gRPC Sever的内部注册中心，这样可以在接收到请求时，通过内部注册中心找到对应的server接口并调用
//3，创建Listener对象，用于监听客户端的连接请求
//4，启动服务器，开始接收客户端请求并调用对应的server接口
