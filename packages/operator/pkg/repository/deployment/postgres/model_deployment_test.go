package postgres_test

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	postgres_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	mdID = "foo"
)

func TestModelDeploymentRepository(t *testing.T) {

	repo := postgres_repo.DeploymentRepo{DB: db}

	as := assert.New(t)
	created := &deployment.ModelDeployment{
		ID: mdID,
		CreatedAt: time.Now().Round(time.Microsecond),
		UpdatedAt: time.Now().Round(time.Microsecond),
		Spec: v1alpha1.ModelDeploymentSpec{
			Image: "not-updated",
		},
	}

	as.NoError(repo.CreateModelDeployment(context.Background(), nil, created))

	as.Exactly(repo.CreateModelDeployment(context.Background(), nil, created), odahuErrors.AlreadyExistError{Entity: mdID})

	fetched, err := repo.GetModelDeployment(context.Background(), nil, mdID)
	as.NoError(err)
	as.Exactly(fetched.ID, created.ID)
	as.Exactly(fetched.Spec, created.Spec)
	as.True(fetched.CreatedAt.Equal(created.CreatedAt))
	t.Logf("fetched createdAt: %+v, initial: %+v", fetched.CreatedAt, created.CreatedAt)
	as.True(fetched.UpdatedAt.Equal(created.UpdatedAt))

	updated := &deployment.ModelDeployment{
		ID: mdID,
		Spec: v1alpha1.ModelDeploymentSpec{
			Image: "updated",
		},
	}

	as.NoError(repo.UpdateModelDeployment(context.Background(), nil, updated))

	fetched, err = repo.GetModelDeployment(context.Background(), nil, mdID)
	as.NoError(err)
	as.Exactly(fetched.Spec, updated.Spec)
	as.Exactly(fetched.Spec.Image, "updated")
	as.True(fetched.UpdatedAt.Equal(updated.UpdatedAt))

	as.NoError(repo.UpdateModelDeploymentStatus(
		context.Background(), nil, mdID, v1alpha1.ModelDeploymentStatus{Replicas: 42}))
	fetched, err = repo.GetModelDeployment(context.Background(), nil, mdID)
	as.NoError(err)
	as.Exactly(fetched.Status.Replicas, int32(42))

	tis, err := repo.GetModelDeploymentList(context.Background(), nil)
	as.NoError(err)
	as.Len(tis, 1)


	as.False(fetched.DeletionMark)
	as.NoError(repo.SetDeletionMark(context.Background(), nil, mdID, true))
	fetched, err = repo.GetModelDeployment(context.Background(), nil, mdID)
	as.NoError(err)
	as.True(fetched.DeletionMark)
	
	as.NoError(repo.DeleteModelDeployment(context.Background(), nil, mdID))
	_, err = repo.GetModelDeployment(context.Background(), nil, mdID)

	as.Error(err)
	as.Exactly(err, odahuErrors.NotFoundError{Entity: mdID})
}
