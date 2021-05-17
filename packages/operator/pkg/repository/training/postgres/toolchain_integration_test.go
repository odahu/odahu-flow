package postgres_test

import (
	_ "github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	postgres_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/postgres"
	. "github.com/onsi/gomega"
	"testing"
)

const (
	tiID            = "foo"
	tiEntrypoint    = "test-entrypoint"
	tiNewEntrypoint = "new-test-entrypoint"
)

func TestToolchainRepository(t *testing.T) {

	tRepo := postgres_repo.ToolchainRepo{DB: db}

	g := NewGomegaWithT(t)

	created := &training.ToolchainIntegration{
		ID: tiID,
		Spec: v1alpha1.ToolchainIntegrationSpec{
			Entrypoint: tiEntrypoint,
		},
	}

	g.Expect(tRepo.SaveToolchainIntegration(created)).NotTo(HaveOccurred())

	g.Expect(tRepo.SaveToolchainIntegration(created)).To(And(
		HaveOccurred(),
		MatchError(odahuErrors.AlreadyExistError{Entity: tiID}),
	))

	fetched, err := tRepo.GetToolchainIntegration(tiID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(created.ID))
	g.Expect(fetched.Spec).To(Equal(created.Spec))

	updated := &training.ToolchainIntegration{
		ID: tiID,
		Spec: v1alpha1.ToolchainIntegrationSpec{
			Entrypoint: tiNewEntrypoint,
		},
	}
	g.Expect(tRepo.UpdateToolchainIntegration(updated)).NotTo(HaveOccurred())

	fetched, err = tRepo.GetToolchainIntegration(tiID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(updated.ID))
	g.Expect(fetched.Spec).To(Equal(updated.Spec))
	g.Expect(fetched.Spec.Entrypoint).To(Equal(tiNewEntrypoint))

	tis, err := tRepo.GetToolchainIntegrationList()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(tis)).To(Equal(1))

	g.Expect(tRepo.DeleteToolchainIntegration(tiID)).NotTo(HaveOccurred())
	_, err = tRepo.GetToolchainIntegration(tiID)
	g.Expect(err).To(And(
		HaveOccurred(),
		MatchError(odahuErrors.NotFoundError{Entity: tiID}),
	))

}
