//
//    Copyright 2019 EPAM Systems
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package kubernetes

import (
	"context"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/filter"
	"time"

	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	logC      = logf.Log.WithName("modelroute--repository")
	MaxSize   = 500
	FirstPage = 0
)

func transform(mr *v1alpha1.ModelRoute) *deployment.ModelRoute {
	return &deployment.ModelRoute{
		ID:     mr.Name,
		Spec:   mr.Spec,
		Status: mr.Status,
	}
}

func (kc *deploymentK8sRepository) GetModelRoute(name string) (*deployment.ModelRoute, error) {
	k8sMR := &v1alpha1.ModelRoute{}
	if err := kc.k8sClient.Get(context.TODO(),
		types.NamespacedName{Name: name, Namespace: kc.namespace},
		k8sMR,
	); err != nil {
		logC.Error(err, "Get Model Route from k8s", "name", name)

		return nil, kubernetes.ConvertK8sErrToOdahuflowErr(err)
	}

	return transform(k8sMR), nil
}

func (kc *deploymentK8sRepository) GetModelRouteList(options ...filter.ListOption) (
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

	labelSelector, err := kubernetes.TransformFilter(listOptions.Filter, "")
	if err != nil {
		logC.Error(err, "Generate label selector")
		return nil, err
	}

	continueToken := ""

	for i := 0; i < *listOptions.Page+1; i++ {
		if err := kc.k8sClient.List(context.TODO(), &k8sMRList, &client.ListOptions{
			LabelSelector: labelSelector,
			Namespace:     kc.namespace,
			Limit:         int64(*listOptions.Size),
			Continue:      continueToken,
		}); err != nil {
			logC.Error(err, "Get Model Route from k8s")

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

func (kc *deploymentK8sRepository) DeleteModelRoute(name string) error {
	conn := &v1alpha1.ModelRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: kc.namespace,
		},
	}

	if err := kc.k8sClient.Delete(context.TODO(),
		conn,
	); err != nil {
		logC.Error(err, "Delete connection from k8s", "name", name)

		return kubernetes.ConvertK8sErrToOdahuflowErr(err)
	}

	return nil
}

func (kc *deploymentK8sRepository) UpdateModelRoute(route *deployment.ModelRoute) error {
	var k8sMR v1alpha1.ModelRoute
	if err := kc.k8sClient.Get(context.TODO(),
		types.NamespacedName{Name: route.ID, Namespace: kc.namespace},
		&k8sMR,
	); err != nil {
		logC.Error(err, "Get route from k8s", "name", route.ID)

		return kubernetes.ConvertK8sErrToOdahuflowErr(err)
	}

	// TODO: think about update, not replacing as for now
	k8sMR.Spec = route.Spec
	k8sMR.Status.UpdatedAt = &metav1.Time{Time: time.Now()}

	if err := kc.k8sClient.Update(context.TODO(), &k8sMR); err != nil {
		logC.Error(err, "Creation of the route", "name", route.ID)

		return kubernetes.ConvertK8sErrToOdahuflowErr(err)
	}

	route.Status = k8sMR.Status

	return nil
}

func (kc *deploymentK8sRepository) CreateModelRoute(route *deployment.ModelRoute) error {
	k8sRoute := &v1alpha1.ModelRoute{
		ObjectMeta: metav1.ObjectMeta{
			Name:      route.ID,
			Namespace: kc.namespace,
		},
		Spec: route.Spec,
	}

	k8sRoute.Status.CreatedAt = &metav1.Time{Time: time.Now()}
	k8sRoute.Status.UpdatedAt = &metav1.Time{Time: time.Now()}

	if err := kc.k8sClient.Create(context.TODO(), k8sRoute); err != nil {
		logC.Error(err, "DataBinding creation error from k8s", "name", route.ID)

		return err
	}

	route.Status = k8sRoute.Status

	return nil
}
