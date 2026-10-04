// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"rpc/gateway/internal/config"
	"rpc/greetrpc/greeterclient"
	"rpc/userrpc/userclient"
	"rpc/videorpc/videoclient"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	Greeter greeterclient.Greeter
	User    userclient.User
	Video   videoclient.Video
}

func NewServiceContext(c config.Config) *ServiceContext {
	gCli := zrpc.MustNewClient(c.Greeter)
	uCli := zrpc.MustNewClient(c.User)
	vCli := zrpc.MustNewClient(c.Video)
	return &ServiceContext{
		Config:  c,
		Greeter: greeterclient.NewGreeter(gCli),
		User:    userclient.NewUser(uCli),
		Video:   videoclient.NewVideo(vCli),
	}
}
