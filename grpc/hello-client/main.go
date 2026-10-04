package main

import (
	"context"
	"fmt"

	pb "grpc/hello-server/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientTokenAuth struct {
}

func (c *ClientTokenAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appId":  "wjh",
		"appKey": "1234567890",
	}, nil
}
func (c *ClientTokenAuth) RequireTransportSecurity() bool {
	return false
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	opts = append(opts, grpc.WithPerRPCCredentials(new(ClientTokenAuth)))

	//连接服务器,第一个参数是服务器地址，第二个参数是传输凭证，这里使用不安全的凭证
	conn, err := grpc.NewClient("localhost:9090", opts...)
	if err != nil {
		fmt.Println("连接失败,err:", err)
		return
	}
	defer conn.Close()

	//建立连接
	client := pb.NewHelloServiceClient(conn)

	//调用服务器接口
	response, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "张三"})
	if err != nil {
		fmt.Println("调用失败,err:", err)
		return
	}
	fmt.Println("调用成功,response:", response.GetReturnmessage())

}
