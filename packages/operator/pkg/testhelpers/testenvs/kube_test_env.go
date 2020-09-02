package testenvs

import (
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	istioschema "github.com/aspenmesh/istio-client-go/pkg/client/clientset/versioned/scheme"
	tektonschema "github.com/tektoncd/pipeline/pkg/client/clientset/versioned/scheme"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

func SetupTestKube(
	crdPaths ...string,
	) (client client.Client, cfg *rest.Config, close func() error, mgr manager.Manager, err error) {

	close = func() error {return nil}

	t := &envtest.Environment{
		// Unit tests can be launched from any of directories because we use the relative path
		// We use the "odahuflow/operator/config/crds" directory here
		CRDDirectoryPaths: crdPaths,
	}
	if err = v1alpha1.AddToScheme(scheme.Scheme); err != nil {
		return client, cfg, close, mgr, err
	}

	istioschema.AddToScheme(scheme.Scheme)

	if err := knservingv1.AddToScheme(scheme.Scheme); err != nil {
		return client, cfg, close, mgr, err
	}
	if err := tektonschema.AddToScheme(scheme.Scheme); err != nil {
		return client, cfg, close, mgr, err
	}

	cfg, err = t.Start()
	if err != nil {
		return client, cfg, close, mgr, err
	}

	close = t.Stop

	mgr, err = manager.New(cfg, manager.Options{NewClient: utils.NewClient, MetricsBindAddress: "0"})
	if err != nil {
		return client, cfg, close, mgr, err
	}
	client = mgr.GetClient()

	return client, cfg, close, mgr, err

}