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

package trainingclient

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	kube_utils "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

const (
	TagKey = "name"
)

var (
	logMT       = logf.Log.WithName("model-training-kube-client")
	MDMaxSize   = 500
	MDFirstPage = 0
	// List of packager steps in execution order
	trainerContainerNames = []string{
		utils.TektonContainerName(odahuflow.TrainerSetupStep),
		utils.TektonContainerName(odahuflow.TrainerTrainStep),
		utils.TektonContainerName(odahuflow.TrainerValidationStep),
		utils.TektonContainerName(odahuflow.TrainerResultStep),
	}
	resultConfigKey = "result"
)

func mtTransformToLabels(mt *training.ModelTraining) map[string]string {
	return map[string]string{
		"toolchain":     mt.Spec.Toolchain,
		"model_name":    mt.Spec.Model.Name,
		"model_version": mt.Spec.Model.Version,
	}
}

func mtTransform(k8sMT *v1alpha1.ModelTraining) *training.ModelTraining {
	return &training.ModelTraining{
		ID:     k8sMT.Name,
		Spec:   k8sMT.Spec,
		Status: k8sMT.Status,
	}
}

type trainingK8SClient struct {
	k8sClient   client.Client
	k8sConfig   *rest.Config
	namespace   string
	tiNamespace string
}

func NewClient(namespace, tiNamespace string, k8sClient client.Client,
	k8sConfig *rest.Config) Client {
	return &trainingK8SClient{
		namespace:   namespace,
		tiNamespace: tiNamespace,
		k8sClient:   k8sClient,
		k8sConfig:   k8sConfig,
	}
}

func (c *trainingK8SClient) SaveModelTrainingResult(id string, result *v1alpha1.TrainingResult) error {
	resultStorage := &corev1.ConfigMap{}
	resultNamespacedName := types.NamespacedName{
		Name:      odahuflow.GenerateTrainingResultCMName(id),
		Namespace: c.namespace,
	}
	if err := c.k8sClient.Get(context.TODO(), resultNamespacedName, resultStorage); err != nil {
		logMT.Error(err, "Result config map must be present", "mt_id", id)
		return err
	}

	oldResult := &v1alpha1.TrainingResult{}
	resultBinary, ok := resultStorage.BinaryData[resultConfigKey]
	if ok {
		if err := json.Unmarshal(resultBinary, oldResult); err != nil {
			return err
		}

		if len(result.CommitID) != 0 {
			oldResult.CommitID = result.CommitID
		}
		if len(result.ArtifactName) != 0 {
			oldResult.ArtifactName = result.ArtifactName
		}
		if len(result.RunID) != 0 {
			oldResult.RunID = result.RunID
		}
	} else {
		oldResult = result
	}

	resultJSON, err := json.Marshal(oldResult)
	if err != nil {
		return err
	}

	resultStorage.BinaryData = map[string][]byte{
		resultConfigKey: resultJSON,
	}

	return c.k8sClient.Update(context.TODO(), resultStorage)
}

func (c *trainingK8SClient) GetModelTrainingResult(id string) (*v1alpha1.TrainingResult, error) {
	resultStorage := &corev1.ConfigMap{}
	resultNamespacedName := types.NamespacedName{
		Name:      odahuflow.GenerateTrainingResultCMName(id),
		Namespace: c.namespace,
	}
	if err := c.k8sClient.Get(context.TODO(), resultNamespacedName, resultStorage); err != nil {
		logMT.Error(err, "Result config map must be present", "mt_id", id)
		return nil, err
	}

	trainingResult := &v1alpha1.TrainingResult{}
	if err := json.Unmarshal(resultStorage.BinaryData[resultConfigKey], &trainingResult); err != nil {
		return nil, err
	}

	return trainingResult, nil
}

func (c *trainingK8SClient) GetModelTraining(id string) (*training.ModelTraining, error) {
	k8sMD := &v1alpha1.ModelTraining{}
	if err := c.k8sClient.Get(context.TODO(),
		types.NamespacedName{Name: id, Namespace: c.namespace},
		k8sMD,
	); err != nil {
		return nil, kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	return mtTransform(k8sMD), nil
}

func (c *trainingK8SClient) GetModelTrainingList(options ...filter.ListOption) (
	[]training.ModelTraining, error,
) {
	var k8sMDList v1alpha1.ModelTrainingList

	listOptions := &filter.ListOptions{
		Filter: nil,
		Page:   &MDFirstPage,
		Size:   &MDMaxSize,
	}
	for _, option := range options {
		option(listOptions)
	}

	labelSelector, err := kube_utils.TransformFilter(listOptions.Filter, TagKey)
	if err != nil {
		logMT.Error(err, "Generate label selector")
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
			logMT.Error(err, "Get Model Training from k8s")

			return nil, err
		}

		continueToken = k8sMDList.ListMeta.Continue
		if *listOptions.Page != i && len(continueToken) == 0 {
			return nil, nil
		}
	}

	mts := make([]training.ModelTraining, len(k8sMDList.Items))
	for i := 0; i < len(k8sMDList.Items); i++ {
		currentMT := k8sMDList.Items[i]

		mts[i] = training.ModelTraining{ID: currentMT.Name, Spec: currentMT.Spec, Status: currentMT.Status}
	}

	return mts, nil
}

func (c *trainingK8SClient) DeleteModelTraining(id string) error {
	mt := &v1alpha1.ModelTraining{
		ObjectMeta: metav1.ObjectMeta{
			Name:      id,
			Namespace: c.namespace,
		},
	}

	if err := c.k8sClient.Delete(context.TODO(), mt); err != nil {
		logMT.Error(err, "Delete Model Training from k8s", "id", id)

		return err
	}

	return nil
}

func (c *trainingK8SClient) UpdateModelTraining(mt *training.ModelTraining) error {
	var k8sMD v1alpha1.ModelTraining
	if err := c.k8sClient.Get(context.TODO(),
		types.NamespacedName{Name: mt.ID, Namespace: c.namespace},
		&k8sMD,
	); err != nil {
		logMT.Error(err, "Get Model Training from k8s", "name", mt.ID)

		return err
	}

	// TODO: think about update, not replacing as for now
	k8sMD.Spec = mt.Spec
	k8sMD.Status.State = v1alpha1.ModelTrainingUnknown
	k8sMD.Status.ExitCode = nil
	k8sMD.Status.Reason = nil
	k8sMD.Status.Message = nil
	k8sMD.ObjectMeta.Labels = mtTransformToLabels(mt)

	if err := c.k8sClient.Update(context.TODO(), &k8sMD); err != nil {
		logMT.Error(err, "Update of the Model Training", "name", mt.ID)

		return kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	mt.Status = k8sMD.Status

	return nil
}

func (c *trainingK8SClient) CreateModelTraining(mt *training.ModelTraining) error {
	k8sMd := &v1alpha1.ModelTraining{
		ObjectMeta: metav1.ObjectMeta{
			Name:      mt.ID,
			Namespace: c.namespace,
			Labels:    mtTransformToLabels(mt),
		},
		Spec: mt.Spec,
	}

	if err := c.k8sClient.Create(context.TODO(), k8sMd); err != nil {
		logMT.Error(err, "StorageEntity creation error from k8s", "name", mt.ID)

		return kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	mt.Status = k8sMd.Status

	return nil
}

func (c *trainingK8SClient) GetModelTrainingLogs(id string, writer utils.Writer, follow bool) error {
	var mt v1alpha1.ModelTraining
	if err := c.k8sClient.Get(context.TODO(),
		types.NamespacedName{Name: id, Namespace: c.namespace},
		&mt,
	); err != nil {
		return err
	} else if mt.Status.State != v1alpha1.ModelTrainingFailed &&
		mt.Status.State != v1alpha1.ModelTrainingRunning &&
		mt.Status.State != v1alpha1.ModelTrainingSucceeded {
		return fmt.Errorf("model Training %s has not started yet", id)
	}

	// Collect logs from all containers in execution order
	for _, containerName := range trainerContainerNames {
		err := utils.StreamLogs(
			c.k8sConfig,
			c.namespace,
			mt.Status.PodName,
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
