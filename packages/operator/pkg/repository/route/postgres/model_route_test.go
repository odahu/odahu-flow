package postgres_test

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	postgres_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/route/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	mrID = "foo"
)

func TestModelRouteRepository(t *testing.T) {

	repo := postgres_repo.RouteRepo{DB: db}

	as := assert.New(t)
	created := &deployment.ModelRoute{
		ID:        mrID,
		CreatedAt: time.Now().Round(time.Microsecond),
		UpdatedAt: time.Now().Round(time.Microsecond),
		Spec: v1alpha1.ModelRouteSpec{
			URLPrefix: "not-updated",
		},
	}

	as.NoError(repo.SaveModelRoute(context.Background(), nil, created))

	as.Exactly(repo.SaveModelRoute(context.Background(), nil, created), odahuErrors.AlreadyExistError{Entity: mrID})

	fetched, err := repo.GetModelRoute(context.Background(), nil, mrID)
	as.NoError(err)
	as.Exactly(fetched.ID, created.ID)
	as.Exactly(fetched.Spec, created.Spec)
	as.True(fetched.CreatedAt.Equal(created.CreatedAt))
	t.Logf("fetched createdAt: %+v, initial: %+v", fetched.CreatedAt, created.CreatedAt)
	as.True(fetched.UpdatedAt.Equal(created.UpdatedAt))

	updated := &deployment.ModelRoute{
		ID: mrID,
		Spec: v1alpha1.ModelRouteSpec{
			URLPrefix: "updated",
		},
	}

	as.NoError(repo.UpdateModelRoute(context.Background(), nil, updated))

	fetched, err = repo.GetModelRoute(context.Background(), nil, mrID)
	as.NoError(err)
	as.Exactly(fetched.Spec, updated.Spec)
	as.Exactly(fetched.Spec.URLPrefix, "updated")
	as.True(fetched.UpdatedAt.Equal(updated.UpdatedAt))

	as.NoError(repo.UpdateModelRouteStatus(
		context.Background(), nil, mrID, v1alpha1.ModelRouteStatus{EdgeURL: "new edge"}))
	fetched, err = repo.GetModelRoute(context.Background(), nil, mrID)
	as.NoError(err)
	as.Exactly(fetched.Status.EdgeURL, "new edge")

	tis, err := repo.GetModelRouteList(context.Background(), nil)
	as.NoError(err)
	as.Len(tis, 1)

	as.False(fetched.DeletionMark)
	as.NoError(repo.SetDeletionMark(context.Background(), nil, mrID, true))
	fetched, err = repo.GetModelRoute(context.Background(), nil, mrID)
	as.NoError(err)
	as.True(fetched.DeletionMark)

	as.NoError(repo.DeleteModelRoute(context.Background(), nil, mrID))
	_, err = repo.GetModelRoute(context.Background(), nil, mrID)

	as.Error(err)
	as.Exactly(err, odahuErrors.NotFoundError{Entity: mrID})
}
