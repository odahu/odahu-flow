package postgres_test

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	odahuErrors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	postgres_repo "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training/postgres"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

const (
	mtID = "foo"
)

type Suite struct {
	suite.Suite
	repo postgres_repo.TrainingRepo
}

func (s *Suite) SetupSuite() {
	s.repo = postgres_repo.TrainingRepo{DB: db}
}

func (s *Suite) TearDownTest() {
	err := s.repo.DeleteModelTraining(context.Background(), nil, mtID)
	if err != nil && !odahuErrors.IsNotFoundError(err) {
		s.T().Fatal(err)
	}
}

func TestRun(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestModelTrainingRepository() {

	tRepo := postgres_repo.TrainingRepo{DB: db}

	g := NewGomegaWithT(s.T())

	created := &training.ModelTraining{
		ID:        mtID,
		CreatedAt: time.Now().Round(time.Microsecond),
		UpdatedAt: time.Now().Round(time.Microsecond),
		Spec: v1alpha1.ModelTrainingSpec{
			WorkDir: "/foo",
		},
	}

	g.Expect(tRepo.SaveModelTraining(context.TODO(), nil, created)).NotTo(HaveOccurred())

	g.Expect(tRepo.SaveModelTraining(context.TODO(), nil, created)).To(And(
		HaveOccurred(),
		MatchError(odahuErrors.AlreadyExistError{Entity: mtID}),
	))

	fetched, err := tRepo.GetModelTraining(context.TODO(), nil, mtID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(created.ID))
	g.Expect(fetched.Spec).To(Equal(created.Spec))
	g.Expect(fetched.CreatedAt.Equal(created.CreatedAt)).Should(BeTrue())
	g.Expect(fetched.UpdatedAt.Equal(created.UpdatedAt)).Should(BeTrue())

	updated := &training.ModelTraining{
		ID:        mtID,
		UpdatedAt: created.CreatedAt.Add(time.Hour),
		Spec: v1alpha1.ModelTrainingSpec{
			WorkDir: "/foo-updated",
		},
	}
	g.Expect(tRepo.UpdateModelTraining(context.TODO(), nil, updated)).NotTo(HaveOccurred())

	fetched, err = tRepo.GetModelTraining(context.TODO(), nil, mtID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(updated.ID))
	g.Expect(fetched.Spec).To(Equal(updated.Spec))
	g.Expect(fetched.Spec.WorkDir).To(Equal("/foo-updated"))
	g.Expect(fetched.UpdatedAt.Equal(updated.UpdatedAt)).Should(BeTrue())

	newStatus := v1alpha1.ModelTrainingStatus{PodName: "Some name"}
	g.Expect(tRepo.UpdateModelTrainingStatus(context.TODO(), nil, mtID, newStatus)).NotTo(HaveOccurred())
	fetched, err = tRepo.GetModelTraining(context.TODO(), nil, mtID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.Status.PodName).To(Equal(newStatus.PodName))

	g.Expect(fetched.DeletionMark).Should(BeFalse())
	g.Expect(tRepo.SetDeletionMark(context.TODO(), nil, mtID, true)).Should(Not(HaveOccurred()))
	fetched, err = tRepo.GetModelTraining(context.TODO(), nil, mtID)
	g.Expect(err).Should(Not(HaveOccurred()))
	g.Expect(fetched.DeletionMark).Should(BeTrue())

	tis, err := tRepo.GetModelTrainingList(context.TODO(), nil)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(len(tis)).To(Equal(1))

	g.Expect(tRepo.DeleteModelTraining(context.TODO(), nil, mtID)).NotTo(HaveOccurred())
	_, err = tRepo.GetModelTraining(context.TODO(), nil, mtID)
	g.Expect(err).To(And(
		HaveOccurred(),
		MatchError(odahuErrors.NotFoundError{Entity: mtID}),
	))

}
