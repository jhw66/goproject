package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"rpc/greetrpc/greeterclient"
	"rpc/userrpc/userclient"
	"rpc/videorpc/videoclient"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
)

// 集成测试：需先启动 greetrpc :8080、userrpc :8081、videorpc :8082，并设置 RUN_RPC_INTEGRATION=1
func TestIntegration_Direct_GreeterUserVideo(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration in short mode")
	}
	if os.Getenv("RUN_RPC_INTEGRATION") == "" {
		t.Skip("set RUN_RPC_INTEGRATION=1 with services on :8080 :8081 :8082")
	}

	root, err := filepath.Abs("../..")
	if err != nil {
		t.Fatal(err)
	}
	cfgPath := filepath.Join(root, "cmd", "client-direct", "etc", "client.yaml")

	var c Config
	conf.MustLoad(cfgPath, &c)

	gCli := zrpc.MustNewClient(c.Greeter)
	uCli := zrpc.MustNewClient(c.User)
	vCli := zrpc.MustNewClient(c.Video)

	g := greeterclient.NewGreeter(gCli)
	u := userclient.NewUser(uCli)
	v := videoclient.NewVideo(vCli)
	ctx := context.Background()

	helloResp, err := g.SayHello(ctx, &greeterclient.SayHelloReq{Name: "integration"})
	if err != nil {
		t.Fatalf("Greeter.SayHello: %v", err)
	}
	if helloResp.Message == "" {
		t.Fatal("empty SayHello message")
	}

	nameResp, err := u.GetName(ctx, &userclient.IdRequest{Id: "1"})
	if err != nil {
		t.Fatalf("User.GetName: %v", err)
	}
	if nameResp.Name == "" {
		t.Fatal("empty GetName name")
	}

	titleResp, err := v.GetTitle(ctx, &videoclient.IdRequest{Id: "1"})
	if err != nil {
		t.Fatalf("Video.GetTitle: %v", err)
	}
	if titleResp.Title == "" {
		t.Fatal("empty GetTitle title")
	}
}
