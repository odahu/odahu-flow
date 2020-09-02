package deploymentclient

import (
	"context"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	kube_utils "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"time"
)

const (
	TagKey = "name"
)


var (
	DefaultMDDeleteOption = metav1.DeletePropagationForeground
	log                   = logf.Log.WithName("modeldeployment-kube-client")
	MaxSize               = 500
	FirstPage             = 0
)

type deployClient struct {
	k8sClient      client.Client
	mdDeleteOption metav1.DeletionPropagation
	namespace      string
}

func NewClient(namespace string, k8sClient client.Client) Client {
	return NewClientWithOptions(namespace, k8sClient, DefaultMDDeleteOption)
}

func NewClientWithOptions(namespace string, k8sClient client.Client,
	mdDeleteOption metav1.DeletionPropagation) Client {
	return &deployClient{
		namespace:      namespace,
		k8sClient:      k8sClient,
		mdDeleteOption: mdDeleteOption,
	}
}


func mdTransformToLabels(md *deployment.ModelDeployment) map[string]string {
	return map[string]string{
		"roleName": *md.Spec.RoleName,
	}
}

func mdTransform(k8sMD *v1alpha1.ModelDeployment) *deployment.ModelDeployment {
	return &deployment.ModelDeployment{
		ID:     k8sMD.Name,
		Spec:   k8sMD.Spec,
		Status: k8sMD.Status,
	}
}

func (c *deployClient) GetModelDeployment(name string) (*deployment.ModelDeployment, error) {
	k8sMD := &v1alpha1.ModelDeployment{}
	if err := c.k8sClient.Get(context.TODO(),
		types.NamespacedName{Name: name, Namespace: c.namespace},
		k8sMD,
	); err != nil {
		log.Error(err, "Get Model Deployment from k8s", "name", name)

		return nil, kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	return mdTransform(k8sMD), nil
}

func (c *deployClient) GetModelDeploymentList(options ...filter.ListOption) (
	[]deployment.ModelDeployment, error,
) {
	var k8sMDList v1alpha1.ModelDeploymentList

	listOptions := &filter.ListOptions{
		Filter: nil,
		Page:   &FirstPage,
		Size:   &MaxSize,
	}

	for _, option := range options {
		option(listOptions)
	}

	labelSelector, err := kube_utils.TransformFilter(listOptions.Filter, TagKey)
	if err != nil {
		log.Error(err, "Generate label selector")
		return nil, err
	}

	continueToken := ""

	for i := 0; i < *listOptions.Page+1; i++ {
		if err := c.k8sClient.List(context.TODO(), &k8sMDList, &client.ListOptions{
			LabelSelector: labelSelector,
			Namespace:     c.namespace,
			Limit:         int64(*listOptions.Size),
			Continue:      continueToken,
		}); err != nil {
			log.Error(err, "Get Model Deployment from k8s")

			return nil, err
		}

		continueToken = k8sMDList.ListMeta.Continue
		if *listOptions.Page != i && len(continueToken) == 0 {
			return nil, nil
		}
	}

	mds := make([]deployment.ModelDeployment, len(k8sMDList.Items))

	for i := 0; i < len(k8sMDList.Items); i++ {
		currentMD := k8sMDList.Items[i]

		mds[i] = deployment.ModelDeployment{ID: currentMD.Name, Spec: currentMD.Spec, Status: currentMD.Status}
	}

	return mds, nil
}

func (c *deployClient) DeleteModelDeployment(name string) error {
	md := &v1alpha1.ModelDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: c.namespace,
		},
	}

	if err := c.k8sClient.Delete(context.TODO(),
		md,
		client.PropagationPolicy(c.mdDeleteOption),
	); err != nil {
		log.Error(err, "Delete Model Deployment from k8s", "name", name)

		return kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	return nil
}

func (c *deployClient) UpdateModelDeployment(md *deployment.ModelDeployment) error {
	var k8sMD v1alpha1.ModelDeployment
	if err := c.k8sClient.Get(context.TODO(),
		types.NamespacedName{Name: md.ID, Namespace: c.namespace},
		&k8sMD,
	); err != nil {
		log.Error(err, "Get Model Deployment from k8s", "name", md.ID)

		return kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	k8sMD.Spec = md.Spec
	k8sMD.Status.UpdatedAt = &metav1.Time{Time: time.Now()}
	k8sMD.ObjectMeta.Labels = mdTransformToLabels(md)

	if err := c.k8sClient.Update(context.TODO(), &k8sMD); err != nil {
		log.Error(err, "Creation of the Model Deployment", "name", md.ID)

		return kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	md.Status = k8sMD.Status

	return nil
}

func (c *deployClient) CreateModelDeployment(md *deployment.ModelDeployment) error {
	k8sMd := &v1alpha1.ModelDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      md.ID,
			Namespace: c.namespace,
			Labels:    mdTransformToLabels(md),
		},
		Spec: md.Spec,
	}

	k8sMd.Status.CreatedAt = &metav1.Time{Time: time.Now()}
	k8sMd.Status.UpdatedAt = &metav1.Time{Time: time.Now()}

	if err := c.k8sClient.Create(context.TODO(), k8sMd); err != nil {
		log.Error(err, "ModelDeployment creation error from k8s", "name", md.ID)

		return err
	}

	md.Status = k8sMd.Status

	return nil
}

func transform(mr *v1alpha1.ModelRoute) *deployment.ModelRoute {
	return &deployment.ModelRoute{
		ID:     mr.Name,
		Spec:   mr.Spec,
		Status: mr.Status,
	}
}

func (c *deployClient) GetModelRoute(name string) (*deployment.ModelRoute, error) {
	k8sMR := &v1alpha1.ModelRoute{}
	if err := c.k8sClient.Get(context.TODO(),
		types.NamespacedName{Name: name, Namespace: c.namespace},
		k8sMR,
	); err != nil {
		log.Error(err, "Get Model Route from k8s", "name", name)

		return nil, kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	return transform(k8sMR), nil
}

func (c *deployClient) GetModelRouteList(options ...filter.ListOption) (
	[]deployment.ModelRoute, error,
) {
	var k8sMRList v1alpha1.ModelRouteList

	listOptions := &filter.ListOptions{
		Filter: nil,
		Page:   &FirstPage,
		Size:   &MaxSize,
	}

	for _, option := range options {
		option(listOptions)
	}

	labelSelector, err := kube_utils.TransformFilter(listOptions.Filter, "")
	if err != nil {
		log.Error(err, "Generate label selector")
		return nil, err
	}

	continueToken := ""

	for i := 0; i < *listOptions.Page+1; i++ {
		if err := c.k8sClient.List(context.TODO(), &k8sMRList, &client.ListOptions{
			LabelSelector: labelSelector,
			Namespace:     c.namespace,
			Limit:         int64(*listOptions.Size),
			Continue:      continueToken,
		}); err != nil {
			log.Error(err, "Get Model Route from k8s")

			return nil, err
		}

		continueToken = k8sMRList.ListMeta.Continue
		if *listOptions.Page != i && len(continueToken) == 0 {
			return nil, nil
		}
	}

	conns := make([]deployment.ModelRoute, len(k8sMRList.Items))
	for i := 0; i < len(k8sMRList.Items); i++ {
		currentMR := k8sMRList.Items[i]

		conns[i] = deployment.ModelRoute{ID: currentMR.Name, Spec: currentMR.Spec, Status: currentMR.Status}
	}

	return conns, nil
}

func (c *deployClient) DeleteModelRoute(name string) error {
	conn := &v1alpha1.ModelRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: c.namespace,
		},
	}

	if err := c.k8sClient.Delete(context.TODO(),
		conn,
	); err != nil {
		log.Error(err, "Delete connection from k8s", "name", name)

		return kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	return nil
}

func (c *deployClient) UpdateModelRoute(route *deployment.ModelRoute) error {
	var k8sMR v1alpha1.ModelRoute
	if err := c.k8sClient.Get(context.TODO(),
		types.NamespacedName{Name: route.ID, Namespace: c.namespace},
		&k8sMR,
	); err != nil {
		log.Error(err, "Get route from k8s", "name", route.ID)

		return kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	// TODO: think about update, not replacing as for now
	k8sMR.Spec = route.Spec
	k8sMR.Status.UpdatedAt = &metav1.Time{Time: time.Now()}

	if err := c.k8sClient.Update(context.TODO(), &k8sMR); err != nil {
		log.Error(err, "Creation of the route", "name", route.ID)

		return kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	route.Status = k8sMR.Status

	return nil
}

func (c *deployClient) CreateModelRoute(route *deployment.ModelRoute) error {
	k8sRoute := &v1alpha1.ModelRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      route.ID,
			Namespace: c.namespace,
		},
		Spec: route.Spec,
	}

	k8sRoute.Status.CreatedAt = &metav1.Time{Time: time.Now()}
	k8sRoute.Status.UpdatedAt = &metav1.Time{Time: time.Now()}

	if err := c.k8sClient.Create(context.TODO(), k8sRoute); err != nil {
		log.Error(err, "DataBinding creation error from k8s", "name", route.ID)

		return err
	}

	route.Status = k8sRoute.Status

	return nil
}

