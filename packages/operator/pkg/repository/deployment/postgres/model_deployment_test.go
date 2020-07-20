package postgres_test

import (
	_ "github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	postgres_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/deployment/postgres"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	mdID = "foo"
	mtID = "foo"
)

func TestModelDeploymentRepository(t *testing.T) {

	repo := postgres_repo.DeploymentPostgresRepo{DB: db}

	as := assert.New(t)

	created := &deployment.ModelDeployment{
		ID: mdID,
		Spec: v1alpha1.ModelDeploymentSpec{
			Image: "not-updated",
		},
	}

	as.NoError(repo.CreateModelDeployment(created))

	as.Exactly(repo.CreateModelDeployment(created), odahuErrors.AlreadyExistError{Entity: mdID})

	fetched, err := repo.GetModelDeployment(mdID)
	as.NoError(err)
	as.Exactly(fetched.ID, created.ID)
	as.Exactly(fetched.Spec, created.Spec)

	updated := &deployment.ModelDeployment{
		ID: mdID,
		Spec: v1alpha1.ModelDeploymentSpec{
			Image: "updated",
		},
	}

	as.NoError(repo.UpdateModelDeployment(updated))

	fetched, err = repo.GetModelDeployment(mdID)
	as.NoError(err)
	as.Exactly(fetched.Spec, updated.Spec)
	as.Exactly(fetched.Spec.Image, "updated")

	tis, err := repo.GetModelDeploymentList()
	as.NoError(err)
	as.Len(tis, 1)

	as.NoError(repo.DeleteModelDeployment(mdID))
	_, err = repo.GetModelDeployment(mdID)

	as.Error(err)
	as.Exactly(err, odahuErrors.NotFoundError{Entity: mdID})
}

func TestModelRouteRepository(t *testing.T) {

	repo := postgres_repo.DeploymentPostgresRepo{DB: db}

	as := assert.New(t)

	created := &deployment.ModelRoute{
		ID: mtID,
		Spec: v1alpha1.ModelRouteSpec{
			URLPrefix: "not-updated",
		},
	}

	as.NoError(repo.CreateModelRoute(created))

	as.Exactly(repo.CreateModelRoute(created), odahuErrors.AlreadyExistError{Entity: mtID})

	fetched, err := repo.GetModelRoute(mtID)
	as.NoError(err)
	as.Exactly(fetched.ID, created.ID)
	as.Exactly(fetched.Spec, created.Spec)

	updated := &deployment.ModelRoute{
		ID: mtID,
		Spec: v1alpha1.ModelRouteSpec{
			URLPrefix: "updated",
		},
	}

	as.NoError(repo.UpdateModelRoute(updated))

	fetched, err = repo.GetModelRoute(mtID)
	as.NoError(err)
	as.Exactly(fetched.Spec, updated.Spec)
	as.Exactly(fetched.Spec.URLPrefix, "updated")

	tis, err := repo.GetModelRouteList()
	as.NoError(err)
	as.Len(tis, 1)

	as.NoError(repo.DeleteModelRoute(mtID))
	_, err = repo.GetModelRoute(mtID)

	as.Error(err)
	as.Exactly(err, odahuErrors.NotFoundError{Entity: mtID})
}