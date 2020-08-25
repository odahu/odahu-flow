package deploymentclient_test

import (
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/testhelpers/testenvs"
	"github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/deploymentclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"path/filepath"
	"testing"
	"log"
	. "github.com/onsi/gomega"
)

const (
	mdImage    = "test:image"
	mdNewImage = "test:new_image"
	mdID       = "test-id"
	testNamespace = "default"
	mrName = "test-mr"
)

var (
	mdRoleName = "test-tole"
	c deploymentclient.Client
)

func TestMain(m *testing.M) {
	os.Exit(WrapperTestMain(m))
}

func WrapperTestMain(m *testing.M) int {
	k8sClient, _, closeF, _, err := testenvs.SetupTestKube(filepath.Join("..", "..", "..", "config", "crds"))
	if err != nil {
		log.Println("Unable to setup kubernetes")
		return -1
	}
	defer func() {
		if err := closeF(); err != nil {
			log.Println("Unable to stop k8s test environment")
		}
	}()
	c = deploymentclient.NewClientWithOptions(testNamespace, k8sClient, metav1.DeletePropagationBackground)
	return m.Run()
}

func TestModelDeploymentsKubeClient(t *testing.T) {
	g := NewGomegaWithT(t)

	created := &deployment.ModelDeployment{
		ID: mdID,
		Spec: v1alpha1.ModelDeploymentSpec{
			Image:    mdImage,
			RoleName: &mdRoleName,
		},
	}

	g.Expect(c.CreateModelDeployment(created)).NotTo(HaveOccurred())

	fetched, err := c.GetModelDeployment(mdID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(created.ID))
	g.Expect(fetched.Spec).To(Equal(created.Spec))

	updated := &deployment.ModelDeployment{
		ID: mdID,
		Spec: v1alpha1.ModelDeploymentSpec{
			Image:    mdNewImage,
			RoleName: &mdRoleName,
		},
	}
	g.Expect(c.UpdateModelDeployment(updated)).NotTo(HaveOccurred())

	fetched, err = c.GetModelDeployment(mdID)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(updated.ID))
	g.Expect(fetched.Spec).To(Equal(updated.Spec))
	g.Expect(fetched.Spec.Image).To(Equal(mdNewImage))

	g.Expect(c.DeleteModelDeployment(mdID)).NotTo(HaveOccurred())
	_, err = c.GetModelDeployment(mdID)
	g.Expect(err).To(HaveOccurred())
}

func TestModelRouteKubeClient(t *testing.T) {
	g := NewGomegaWithT(t)

	urlPrefixValue := "/test"
	newURLPrefixValue := "/new/test"
	created := &deployment.ModelRoute{
		ID: mrName,
		Spec: v1alpha1.ModelRouteSpec{
			URLPrefix: urlPrefixValue,
			ModelDeploymentTargets: []v1alpha1.ModelDeploymentTarget{
				{
					Name: mdID,
				},
			},
		},
	}

	g.Expect(c.CreateModelRoute(created)).NotTo(HaveOccurred())

	fetched, err := c.GetModelRoute(mrName)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(created.ID))
	g.Expect(fetched.Spec).To(Equal(created.Spec))

	updated := &deployment.ModelRoute{
		ID: mrName,
		Spec: v1alpha1.ModelRouteSpec{
			URLPrefix: urlPrefixValue,
			ModelDeploymentTargets: []v1alpha1.ModelDeploymentTarget{
				{
					Name: mdID,
				},
			},
		},
	}
	updated.Spec.URLPrefix = newURLPrefixValue
	g.Expect(c.UpdateModelRoute(updated)).NotTo(HaveOccurred())

	fetched, err = c.GetModelRoute(mrName)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(fetched.ID).To(Equal(updated.ID))
	g.Expect(fetched.Spec).To(Equal(updated.Spec))
	g.Expect(fetched.Spec.URLPrefix).To(Equal(newURLPrefixValue))

	g.Expect(c.DeleteModelRoute(mrName)).NotTo(HaveOccurred())
	_, err = c.GetModelRoute(mrName)
	g.Expect(err).To(HaveOccurred())
}