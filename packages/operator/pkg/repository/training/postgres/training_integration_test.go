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

func TestTrainingIntegrationRepository(t *testing.T) {

	tRepo := postgres_repo.TrainingIntegrationRepo{DB: db}

	g := NewGomegaWithT(t)

	created := &training.TrainingIntegration{
		ID: tiID,
		Spec: v1alpha1.TrainingIntegrationSpec{
			Entrypoint: tiEntrypoint,
		},
	}

	g.Expect(tRepo.SaveTrainingIntegration(created)).NotTo(HaveOccurred())

	g.Expect(tRepo.SaveTrainingIntegration(created)).To(And(
		HaveOccurred(),
		MatchError(odahuErrors.AlreadyExistError{Entity: tiID}),
	))

	fetched, err := tRepo.GetTrainingIntegration(tiID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(created.ID))
	g.Expect(fetched.Spec).To(Equal(created.Spec))

	updated := &training.TrainingIntegration{
		ID: tiID,
		Spec: v1alpha1.TrainingIntegrationSpec{
			Entrypoint: tiNewEntrypoint,
		},
	}
	g.Expect(tRepo.UpdateTrainingIntegration(updated)).NotTo(HaveOccurred())

	fetched, err = tRepo.GetTrainingIntegration(tiID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(updated.ID))
	g.Expect(fetched.Spec).To(Equal(updated.Spec))
	g.Expect(fetched.Spec.Entrypoint).To(Equal(tiNewEntrypoint))

	tis, err := tRepo.GetTrainingIntegrationList()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(tis)).To(Equal(1))

	g.Expect(tRepo.DeleteTrainingIntegration(tiID)).NotTo(HaveOccurred())
	_, err = tRepo.GetTrainingIntegration(tiID)
	g.Expect(err).To(And(
		HaveOccurred(),
		MatchError(odahuErrors.NotFoundError{Entity: tiID}),
	))

}
