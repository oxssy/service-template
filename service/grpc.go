package service

import (
	"io/ioutil"
	"net"
	"os"

	template "github.com/oxssy/service-template"
	"github.com/oxssy/service-template/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

type GRPCImplementation interface {
	template.Implementation
	Register(srv *grpc.Server)
}

type gRPCService struct {
	GRPCImplementation
	grpcConfig *config.NetConfig
	listener   net.Listener
}

func (svc *gRPCService) OnReady() error {
	err := svc.GRPCImplementation.OnReady()
	if err != nil {
		return err
	}
	grpcs := grpc.NewServer()
	svc.Register(grpcs)
	reflection.Register(grpcs)
	listener, err := svc.grpcConfig.Listen()
	if err != nil {
		return err
	}
	svc.listener = listener
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(ioutil.Discard, ioutil.Discard, os.Stdout))
	return grpcs.Serve(listener)
}

func (svc *gRPCService) OnClose() error {
	svc.listener.Close()
	return svc.GRPCImplementation.OnClose()
}

func NewGRPCService(impl GRPCImplementation) *template.Service {
	grpcConf := &config.NetConfig{}
	srv := template.NewService(&gRPCService{
		GRPCImplementation: impl,
		grpcConfig:         grpcConf,
	})
	srv.GetConfig().Set("grpc", grpcConf)
	return srv
}
