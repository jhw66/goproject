package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"rpc/greetrpc/greeterclient"
	"rpc/userrpc/userclient"
	"rpc/videorpc/videoclient"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	Greeter zrpc.RpcClientConf
	User    zrpc.RpcClientConf
	Video   zrpc.RpcClientConf
}

func main() {
	configFile := flag.String("f", "cmd/client-direct/etc/client.yaml", "配置文件路径（相对模块根目录 rpc）")
	name := flag.String("name", "go-zero", "SayHello 参数 name")
	id := flag.String("id", "1", "GetName 参数 id")
	vid := flag.String("vid", "1", "GetTitle 参数 id")
	flag.Parse()

	var c Config
	conf.MustLoad(*configFile, &c)

	gCli := zrpc.MustNewClient(c.Greeter)
	uCli := zrpc.MustNewClient(c.User)
	vCli := zrpc.MustNewClient(c.Video)

	g := greeterclient.NewGreeter(gCli)
	u := userclient.NewUser(uCli)
	v := videoclient.NewVideo(vCli)

	ctx := context.Background()

	helloResp, err := g.SayHello(ctx, &greeterclient.SayHelloReq{Name: *name})
	if err != nil {
		log.Fatalf("Greeter.SayHello: %v", err)
	}
	fmt.Printf("Greeter: %s\n", helloResp.Message)

	nameResp, err := u.GetName(ctx, &userclient.IdRequest{Id: *id})
	if err != nil {
		log.Fatalf("User.GetName: %v", err)
	}
	fmt.Printf("User.GetName: %s\n", nameResp.Name)

	titleResp, err := v.GetTitle(ctx, &videoclient.IdRequest{Id: *vid})
	if err != nil {
		log.Fatalf("Video.GetTitle: %v", err)
	}
	fmt.Printf("Video.GetTitle: %s\n", titleResp.Title)
}
