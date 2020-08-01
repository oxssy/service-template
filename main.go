package template

import (
	"io/ioutil"
	"net/http"
	"os"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
)

// StartGRPCService runs a given GRPCService.
func StartGRPCService(svc GRPCService) error {
	config, conn := makeConnection()
	grpcConf := &NetConfig{}
	config.Set("grpc", grpcConf)
	err := initalize(svc, config)
	if err != nil {
		return err
	}
	err = connect(svc, conn)
	if err != nil {
		return err
	}
	defer close(svc, conn)
	err = svc.OnReady()
	if err != nil {
		return err
	}
	grpcs := grpc.NewServer()
	svc.Register(grpcs)
	reflection.Register(grpcs)
	listener, err := grpcConf.Listen()
	if err != nil {
		return err
	}
	defer listener.Close()
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(ioutil.Discard, ioutil.Discard, os.Stdout))
	return grpcs.Serve(listener)
}

// StartHTTPService runs a given HTTPService.
func StartHTTPService(svc HTTPService) error {
	config, conn := makeConnection()
	httpConf := &NetConfig{}
	config.Set("http", httpConf)
	err := initalize(svc, config)
	if err != nil {
		return err
	}
	err = connect(svc, conn)
	if err != nil {
		return err
	}
	defer close(svc, conn)
	err = svc.OnReady()
	if err != nil {
		return err
	}
	router := mux.NewRouter()
	svc.Route(router)
	listener, err := httpConf.Listen()
	if err != nil {
		return err
	}
	defer listener.Close()
	return http.Serve(listener, router)
}

func makeConnection() (ConfigMap, *Connection) {
	config := NewConfigMap()
	return config, NewConnection(config)
}

func initalize(svc Service, config ConfigMap) error {
	svc.Setup(config)
	err := config.Load()
	if err != nil {
		return err
	}
	return svc.OnConfig(config)
}

func connect(svc Service, conn *Connection) error {
	err := conn.Connect()
	if err != nil {
		return err
	}
	return svc.OnConnect(conn)
}

func close(svc Service, conn *Connection) {
	err := svc.OnClose()
	if err != nil {
		glog.Error(err)
	}
	err = conn.Close()
	if err != nil {
		glog.Error(err)
	}
}
