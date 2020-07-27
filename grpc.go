package template

import (
	"io/ioutil"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

type GRPCImplementation interface {
	Implementation
	Register(srv *grpc.Server)
}

type gRPCService struct {
	GRPCImplementation
	grpcConfig *NetConfig
	listener   net.Listener
}

func (svc *gRPCService) Setup(cm ConfigMap) {
	svc.grpcConfig = &NetConfig{}
	cm.Set("grpc", svc.grpcConfig)
	svc.GRPCImplementation.Setup(cm)
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

func NewGRPCService(impl GRPCImplementation) *Service {
	return NewService(&gRPCService{GRPCImplementation: impl})
}
