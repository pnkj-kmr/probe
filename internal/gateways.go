package internal

import (
	"context"
	"io"
)

// Gateway instance represent the external system connectivity
// here gateway means calling external system from probe
type Gateway interface {
	io.Writer
	io.Closer

	Open() error
	OpenWithContext(context.Context) error
}

// GatewayStream instance represent the external system connectivity
// with long running connection like kafka, rabbitmq etc
type GatewayStream interface {
	io.Writer
	io.Closer

	Open() error
	OpenWithContext(context.Context) error
}
