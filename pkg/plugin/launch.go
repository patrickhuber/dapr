package plugin

import (
	"context"
	"net"
)

type ComponentType string

const (
	StateComponentType ComponentType = "state"
)

type LaunchContext struct {
	Address net.Addr
}

type AsyncLauncher interface {
	Launch(ctx context.Context, metadata *Metadata) error
}
