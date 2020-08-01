package template

import (
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

// Service contains all method a service must svcement.
type Service interface {
	// Setup is called at the beginning of the service life cycle.
	// Typically, your service Service should add necessary configurations to the ConfigMap.
	Setup(config ConfigMap)

	// OnConfig is called when the ConfigMap is loaded from the environment.
	// Typically, your service Service can use the loaded configurations to set itself up.
	OnConfig(config ConfigMap) error

	// OnConnect is called when the Connection object makes all its connections to data stores.
	// Typically, your service Service will save various connections to its own struct.
	OnConnect(conn *Connection) error

	// OnReady is called when the service is ready to serve.
	// Your serivce Service should use this method to start listening to requests.
	OnReady() error

	// OnClose is called when the service is about to close.
	// Your service Service should use this method to cleanup before the Connection object
	// severs all its connections.
	OnClose() error
}

// GRPCService is a Service that has gRPC endpoints.
type GRPCService interface {
	Service
	// Register the GRPCService to a gRPC Server.
	Register(srv *grpc.Server)
}

// HTTPService is a Service that has HTTP endpoints.
type HTTPService interface {
	Service
	// Route sets up handlers for the HTTP endpoints in a Router.
	Route(router *mux.Router)
}
