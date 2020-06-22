package template

// Implementation contains all method a service must implement.
type Implementation interface {
	// Setup is called at the beginning of the service life cycle.
	// Typically, your service Implementation should add necessary configurations to the ConfigMap.
	Setup(config ConfigMap)

	// OnConfig is called when the ConfigMap is loaded from the environment.
	// Typically, your service Implementation can use the loaded configurations to set itself up.
	OnConfig(config ConfigMap) error

	// OnConnect is called when the Connection object makes all its connections to data stores.
	// Typically, your service Implementation will save various connections to its own struct.
	OnConnect(conn *Connection) error

	// OnReady is called when the service is ready to serve.
	// Your serivce Implementation should use this method to start listening to requests.
	OnReady() error

	// OnClose is called when the service is about to close.
	// Your service Implementation should use this method to cleanup before the Connection object
	// severs all its connections.
	OnClose() error
}

// Service represents a HTTP or gRPC service.
type Service struct {
	Implementation
	config     ConfigMap
	connection *Connection
}

// NewService creates a service object from a service Implementation.
func NewService(impl Implementation) *Service {
	config := NewConfigMap()
	return &Service{
		Implementation: impl,
		config:         config,
		connection:     NewConnection(config),
	}
}

// GetConfig returns the ConfigMap of this service.
func (svc *Service) GetConfig() ConfigMap {
	return svc.config
}

// GetConnection returns the Connection of this service.
func (svc *Service) GetConnection() *Connection {
	return svc.connection
}

// Start the service.
// Calling this will load config from the environment, make data store connections,
// and finally start listening to requests.
// Your service Implementation's various life cycle hooks will be called during this process.
func (svc *Service) Start() error {
	err := svc.initalize()
	if err != nil {
		return err
	}
	err = svc.connect()
	if err != nil {
		return err
	}
	return svc.OnReady()
}

// Close will close all connections from this service,
// after any cleanup in the Implementation's OnClose method.
func (svc *Service) Close() error {
	err := svc.OnClose()
	if connErr := svc.connection.Close(); err == nil && connErr != nil {
		err = connErr
	}
	return err
}

func (svc *Service) initalize() error {
	svc.Setup((svc.config))
	err := svc.config.Load()
	if err != nil {
		return err
	}
	return svc.OnConfig(svc.config)
}

func (svc *Service) connect() error {
	err := svc.connection.Connect()
	if err != nil {
		return err
	}
	return svc.OnConnect(svc.connection)
}
