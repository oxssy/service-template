package service

import (
	"net"
	"net/http"

	template "github.com/oxssy/service-template"
	"github.com/oxssy/service-template/config"

	"github.com/gorilla/mux"
)

type HTTPImplementation interface {
	template.Implementation
	Routes(router *mux.Router)
}

type httpService struct {
	HTTPImplementation
	httpConfig *config.NetConfig
	listener   net.Listener
}

func (svc *httpService) OnReady() error {
	err := svc.HTTPImplementation.OnReady()
	if err != nil {
		return err
	}
	router := mux.NewRouter()
	svc.Routes(router)
	listener, err := svc.httpConfig.Listen()
	if err != nil {
		return err
	}
	svc.listener = listener
	return http.Serve(svc.listener, router)
}

func NewHTTPService(impl HTTPImplementation) *template.Service {
	httpConf := &config.NetConfig{}
	srv := template.NewService(&httpService{
		HTTPImplementation: impl,
		httpConfig:         httpConf,
	})
	srv.GetConfig().Set("http", httpConf)
	return srv
}
