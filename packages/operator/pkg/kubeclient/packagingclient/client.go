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

package packagingclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"k8s.io/client-go/rest"
	"time"

	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	kube_utils "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	TagKey = "name"
)


var (
	logMP       = logf.Log.WithName("model-packaging-kube-client")
	MpMaxSize   = 500
	MpFirstPage = 0
	// List of packager steps in execution order
	packagerContainerNames = []string{
		utils.TektonContainerName(odahuflow.PackagerSetupStep),
		utils.TektonContainerName(odahuflow.PackagerPackageStep),
		utils.TektonContainerName(odahuflow.PackagerResultStep),
	}
	resultConfigKey = "result"
)

func TransformMpFromK8s(mp *v1alpha1.ModelPackaging) (*packaging.ModelPackaging, error) {
	var arguments map[string]interface{}
	if err := json.Unmarshal([]byte(mp.Spec.Arguments), &arguments); err != nil {
		return nil, err
	}

	return &packaging.ModelPackaging{
		ID: mp.Name,
		Spec: packaging.ModelPackagingSpec{
			ArtifactName:     *mp.Spec.ArtifactName,
			IntegrationName:  mp.Spec.Type,
			Image:            mp.Spec.Image,
			Arguments:        arguments,
			Targets:          mp.Spec.Targets,
			Resources:        mp.Spec.Resources,
			OutputConnection: mp.Spec.OutputConnection,
		},
		Status: mp.Status,
	}, nil
}

// Take a look to the documentation of TransformPackagingIntegrationFromK8s
func TransformMpToK8s(mp *packaging.ModelPackaging, k8sNamespace string) (*v1alpha1.ModelPackaging, error) {

	argumentsBytes, err := json.Marshal(mp.Spec.Arguments)
	if err != nil {
		return nil, err
	}

	return &v1alpha1.ModelPackaging{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mp.ID,
			Namespace: k8sNamespace,
			Labels: map[string]string{
				"type": mp.Spec.IntegrationName,
			},
		},
		Spec: v1alpha1.ModelPackagingSpec{
			ArtifactName:     &mp.Spec.ArtifactName,
			Type:             mp.Spec.IntegrationName,
			Image:            mp.Spec.Image,
			Arguments:        string(argumentsBytes),
			Targets:          mp.Spec.Targets,
			Resources:        mp.Spec.Resources,
			OutputConnection: mp.Spec.OutputConnection,
		},
	}, nil
}

type packagingK8SClient struct {
	k8sClient   client.Client
	k8sConfig   *rest.Config
	namespace   string
	piNamespace string
}

func NewClient(namespace, piNamespace string, k8sClient client.Client,
	k8sConfig *rest.Config) Client {
	return &packagingK8SClient{
		namespace:   namespace,
		k8sClient:   k8sClient,
		piNamespace: piNamespace,
		k8sConfig:   k8sConfig,
	}
}


func (c *packagingK8SClient) GetModelPackaging(id string) (*packaging.ModelPackaging, error) {
	k8sMp := &v1alpha1.ModelPackaging{}
	if err := c.k8sClient.Get(context.TODO(),
		types.NamespacedName{Name: id, Namespace: c.namespace},
		k8sMp,
	); err != nil {
		logMP.Error(err, "Get Model Packaging from k8s", "id", id)

		return nil, kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	return TransformMpFromK8s(k8sMp)
}

func (c *packagingK8SClient) GetModelPackagingList(options ...filter.ListOption) (
	[]packaging.ModelPackaging, error,
) {
	var k8sMpList v1alpha1.ModelPackagingList

	listOptions := &filter.ListOptions{
		Filter: nil,
		Page:   &MpFirstPage,
		Size:   &MpMaxSize,
	}
	for _, option := range options {
		option(listOptions)
	}

	labelSelector, err := kube_utils.TransformFilter(listOptions.Filter, TagKey)
	if err != nil {
		logMP.Error(err, "Generate label selector")
		return nil, err
	}
	continueToken := ""

	for i := 0; i < *listOptions.Page+1; i++ {
		if err := c.k8sClient.List(context.TODO(), &k8sMpList, &client.ListOptions{
			LabelSelector: labelSelector,
			Namespace:     c.namespace,
			Limit:         int64(*listOptions.Size),
			Continue:      continueToken,
		}); err != nil {
			logMP.Error(err, "Get Model Packaging from k8s")

			return nil, err
		}

		continueToken = k8sMpList.ListMeta.Continue
		if *listOptions.Page != i && len(continueToken) == 0 {
			return nil, nil
		}
	}

	mps := make([]packaging.ModelPackaging, len(k8sMpList.Items))
	for i := 0; i < len(k8sMpList.Items); i++ {
		k8sMp := k8sMpList.Items[i]

		mp, err := TransformMpFromK8s(&k8sMp)
		if err != nil {
			return nil, err
		}

		mps[i] = *mp
	}

	return mps, nil
}

func (c *packagingK8SClient) DeleteModelPackaging(id string) error {
	mp := &v1alpha1.ModelPackaging{
		ObjectMeta: metav1.ObjectMeta{
			Name:      id,
			Namespace: c.namespace,
		},
	}

	if err := c.k8sClient.Delete(context.TODO(), mp); err != nil {
		logMP.Error(err, "Delete Model Packaging from k8s", "id", id)

		return kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	return nil
}

func (c *packagingK8SClient) UpdateModelPackaging(mp *packaging.ModelPackaging) error {
	var k8sMp v1alpha1.ModelPackaging
	if err := c.k8sClient.Get(context.TODO(),
		types.NamespacedName{Name: mp.ID, Namespace: c.namespace},
		&k8sMp,
	); err != nil {
		logMP.Error(err, "Get Model Packaging from k8s", "id", mp.ID)

		return kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	// TODO: think about update, not replacing as for now
	updatedK8sMpSpec, err := TransformMpToK8s(mp, c.namespace)
	if err != nil {
		return err
	}

	k8sMp.Spec = updatedK8sMpSpec.Spec
	k8sMp.Status.State = v1alpha1.ModelPackagingUnknown
	k8sMp.Status.ExitCode = nil
	k8sMp.Status.Reason = nil
	k8sMp.Status.Message = nil
	k8sMp.Status.Results = []v1alpha1.ModelPackagingResult{}
	k8sMp.Status.UpdatedAt = &metav1.Time{Time: time.Now()}
	k8sMp.ObjectMeta.Labels = updatedK8sMpSpec.Labels

	if err := c.k8sClient.Update(context.TODO(), &k8sMp); err != nil {
		logMP.Error(err, "Creation of the Model Packaging", "id", mp.ID)

		return err
	}

	mp.Status = k8sMp.Status

	return nil
}

func (c *packagingK8SClient) CreateModelPackaging(mp *packaging.ModelPackaging) error {
	k8sMp, err := TransformMpToK8s(mp, c.namespace)
	if err != nil {
		return err
	}

	k8sMp.Status.CreatedAt = &metav1.Time{Time: time.Now()}
	k8sMp.Status.UpdatedAt = &metav1.Time{Time: time.Now()}

	if err := c.k8sClient.Create(context.TODO(), k8sMp); err != nil {
		logMP.Error(err, "Model packaging creation error from k8s", "id", mp.ID)

		return err
	}

	mp.Status = k8sMp.Status

	return nil
}

func (c *packagingK8SClient) GetModelPackagingLogs(id string, writer utils.Writer, follow bool) error {
	var mp v1alpha1.ModelPackaging
	if err := c.k8sClient.Get(context.TODO(),
		types.NamespacedName{Name: id, Namespace: c.namespace},
		&mp,
	); err != nil {
		return err
	} else if mp.Status.State != v1alpha1.ModelPackagingSucceeded &&
		mp.Status.State != v1alpha1.ModelPackagingRunning &&
		mp.Status.State != v1alpha1.ModelPackagingFailed {
		return fmt.Errorf("model packaing %s has not started yet", id)
	}

	// Collect logs from all containers in execution order
	for _, containerName := range packagerContainerNames {
		err := utils.StreamLogs(
			c.k8sConfig,
			c.namespace,
			mp.Status.PodName,
			containerName,
			writer,
			follow,
			utils.LogFlushSize,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *packagingK8SClient) SaveModelPackagingResult(id string, result []v1alpha1.ModelPackagingResult) error {
	resultStorage := &corev1.ConfigMap{}
	resultNamespacedName := types.NamespacedName{
		Name:      odahuflow.GeneratePackageResultCMName(id),
		Namespace: c.namespace,
	}
	if err := c.k8sClient.Get(context.TODO(), resultNamespacedName, resultStorage); err != nil {
		logMP.Error(err, "Result config map must be present", "mp_id", id)
		return err
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resultStorage.BinaryData = map[string][]byte{
		resultConfigKey: resultJSON,
	}

	return c.k8sClient.Update(context.TODO(), resultStorage)
}

func (c *packagingK8SClient) GetModelPackagingResult(id string) ([]v1alpha1.ModelPackagingResult, error) {
	resultStorage := &corev1.ConfigMap{}
	resultNamespacedName := types.NamespacedName{
		Name:      odahuflow.GeneratePackageResultCMName(id),
		Namespace: c.namespace,
	}
	if err := c.k8sClient.Get(context.TODO(), resultNamespacedName, resultStorage); err != nil {
		logMP.Error(err, "Result config map must be present", "mp_id", id)
		return nil, err
	}

	packResult := make([]v1alpha1.ModelPackagingResult, 0)
	if err := json.Unmarshal(resultStorage.BinaryData[resultConfigKey], &packResult); err != nil {
		return nil, err
	}

	return packResult, nil
}
