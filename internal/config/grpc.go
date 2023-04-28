package config

import (
	"fmt"
	"strings"
)

// GrpcConfig for application
type GrpcConfig struct {
	// BindRaw port string, default "8080"
	BindRaw string `envconfig:"API_GRPC_PORT" default:"8080"`

	// ----------------------------
	// Calculated config parameters
	Bind string
}

func (c *GrpcConfig) GetBindPort() string {
	return c.Bind
}

// Prepare variables to static configuration
func (c *GrpcConfig) Prepare() error {
	c.Bind = fmt.Sprintf(":%s", strings.TrimLeft(c.BindRaw, ":"))

	return nil
}

// PrepareWith variables with dependencies service-components
func (c *GrpcConfig) PrepareWith(cfgSrvList ...interface{}) error {
	return nil
}
