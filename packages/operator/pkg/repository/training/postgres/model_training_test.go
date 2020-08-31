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
	mtID = "foo"
)

func TestModelTrainingRepository(t *testing.T) {

	tRepo := postgres_repo.TrainingRepo{DB: db}

	g := NewGomegaWithT(t)

	created := &training.ModelTraining{
		ID: mtID,
		Spec: v1alpha1.ModelTrainingSpec{
			WorkDir: "/foo",
		},
	}

	g.Expect(tRepo.CreateModelTraining(created)).NotTo(HaveOccurred())

	g.Expect(tRepo.CreateModelTraining(created)).To(And(
		HaveOccurred(),
		MatchError(odahuErrors.AlreadyExistError{Entity: mtID}),
	))

	fetched, err := tRepo.GetModelTraining(mtID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(created.ID))
	g.Expect(fetched.Spec).To(Equal(created.Spec))

	updated := &training.ModelTraining{
		ID: mtID,
		Spec: v1alpha1.ModelTrainingSpec{
			WorkDir: "/foo-updated",
		},
	}
	g.Expect(tRepo.UpdateModelTraining(updated)).NotTo(HaveOccurred())

	fetched, err = tRepo.GetModelTraining(mtID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(updated.ID))
	g.Expect(fetched.Spec).To(Equal(updated.Spec))
	g.Expect(fetched.Spec.WorkDir).To(Equal("/foo-updated"))

	newStatus := v1alpha1.ModelTrainingStatus{PodName: "Some name"}
	g.Expect(tRepo.UpdateModelTrainingStatus(mtID, newStatus)).NotTo(HaveOccurred())
	fetched, err = tRepo.GetModelTraining(mtID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.Status.PodName).To(Equal(newStatus.PodName))

	g.Expect(fetched.DeletionMark).Should(BeFalse())
	g.Expect(tRepo.SetDeletionMark(mtID, true)).Should(Not(HaveOccurred()))
	fetched, err = tRepo.GetModelTraining(mtID)
	g.Expect(err).Should(Not(HaveOccurred()))
	g.Expect(fetched.DeletionMark).Should(BeTrue())

	tis, err := tRepo.GetModelTrainingList()
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(tis)).To(Equal(1))

	g.Expect(tRepo.DeleteModelTraining(mtID)).NotTo(HaveOccurred())
	_, err = tRepo.GetModelTraining(mtID)
	g.Expect(err).To(And(
		HaveOccurred(),
		MatchError(odahuErrors.NotFoundError{Entity: mtID}),
	))

}
