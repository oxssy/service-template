package template

import (
	"github.com/golang/glog"
)

func RunService(svc *Service) {
	defer close(svc)
	glog.Error(svc.Start().Error())
}

func close(svc *Service) {
	err := svc.Close()
	if err != nil {
		glog.Error(err.Error())
	}
}
