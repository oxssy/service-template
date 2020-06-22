package config

import (
	"fmt"
	"net"

	template "github.com/oxssy/service-template"
)

// NetConfig contains the host and port parameters for a network listener.
type NetConfig struct {
	Host     string `default:"0.0.0.0"`
	Port     string `default:"80"`
	Protocol string `default:"tcp"`
}

// ConfigType of NetConfig is NET.
func (c *NetConfig) ConfigType() template.ConfigTypeValue {
	return template.ConfigType.Net
}

// Address returns a string address with host and port.
func (c *NetConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// Listen returns a net.Listener as specified by this NetConfig.
func (c *NetConfig) Listen() (net.Listener, error) {
	return net.Listen(c.Protocol, c.Address())
}
