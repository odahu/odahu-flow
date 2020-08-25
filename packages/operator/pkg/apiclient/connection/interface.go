package connection

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
)

type Client interface {
	GetConnection(id string) (*connection.Connection, error)
}
