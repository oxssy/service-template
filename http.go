package template

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPImplementation interface {
	Implementation
	Routes(router *mux.Router)
}

type httpService struct {
	HTTPImplementation
	httpConfig *NetConfig
	listener   net.Listener
}

func (svc *httpService) Setup(cm ConfigMap) {
	svc.httpConfig = &NetConfig{}
	cm.Set("http", svc.httpConfig)
	svc.HTTPImplementation.Setup(cm)
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

func NewHTTPService(impl HTTPImplementation) *Service {
	return NewService(&httpService{HTTPImplementation: impl})
}
