package route

import (
	"context"
	"database/sql"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
)

type Repository interface {
	GetModelRoute(name string) (*deployment.ModelRoute, error)
	GetModelRouteList(options ...filter.ListOption) ([]deployment.ModelRoute, error)
	DeleteModelRoute(name string) error
	UpdateModelRoute(md *deployment.ModelRoute) error
	CreateModelRoute(md *deployment.ModelRoute) error
	SetDeletionMark(ctx context.Context, tx *sql.Tx, id string, value bool) error
	BeginTransaction(ctx context.Context) (*sql.Tx, error)
}