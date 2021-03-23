/*
 *
 *     Copyright 2021 EPAM Systems
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 */

package batchclient

import (
	"context"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/controllers/utils/batchinferenceutils"
	kube_utils "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log  = logf.Log.WithName("batch-job-kube-client")

var (
	jobContainerNames = []string{
		utils.TektonContainerName(batchinferenceutils.StepSyncData),
		utils.TektonContainerName(batchinferenceutils.StepSyncModel),
		utils.TektonContainerName(batchinferenceutils.StepValidateInput),
		utils.TektonContainerName(batchinferenceutils.StepLogInput),
		utils.TektonContainerName(batchinferenceutils.StepUserContainer),
		utils.TektonContainerName(batchinferenceutils.StepValidateOutput),
		utils.TektonContainerName(batchinferenceutils.StepLogOutput),
		utils.TektonContainerName(batchinferenceutils.StepSyncOutput),
	}
)

type Client struct {
	client client.Client
	namespace string
	cfg *rest.Config
}

func NewClient(client client.Client, namespace string, cfg *rest.Config) *Client {
	return &Client{client: client, namespace: namespace, cfg: cfg}
}

func (c *Client) DeleteInferenceJob(name string) error {
	job := &v1alpha1.BatchInferenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: c.namespace,
		},
	}

	if err := c.client.Delete(context.TODO(), job); err != nil {
		log.Error(err, "Delete Batch Inference from k8s", "id", name)

		return kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}

	return nil
}

func (c *Client) GetInferenceJob(name string) (*v1alpha1.BatchInferenceJob, error) {
	job := &v1alpha1.BatchInferenceJob{}
	if err := c.client.Get(context.TODO(),
		types.NamespacedName{Name: name, Namespace: c.namespace},
		job,
	); err != nil {
		return nil, kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}
	return job, nil
}

func (c *Client) CreateInferenceJob(name string, spec *v1alpha1.BatchInferenceJobSpec) error {
	job := &v1alpha1.BatchInferenceJob{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Namespace: c.namespace,
		},
		Spec:       *spec,
	}
	if err := c.client.Create(context.TODO(), job); err != nil {
		log.Error(err, "Unable to create inference job in kubernetes")
		return kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}
	return nil
}

func (c *Client) ListInferenceJob() (v1alpha1.BatchInferenceJobList, error) {
	var result v1alpha1.BatchInferenceJobList

	if err := c.client.List(context.TODO(), &result); err != nil {
		log.Error(err,"Unable to list inference jobs")
		return result, kube_utils.ConvertK8sErrToOdahuflowErr(err)
	}
	return result, nil
}

func (c *Client) LogInferenceJob(id string, writer utils.Writer, follow bool) error {
	var job v1alpha1.BatchInferenceJob
	if err := c.client.Get(context.TODO(),
		types.NamespacedName{Name: id, Namespace: c.namespace},
		&job,
	); err != nil {
		return err
	} else if job.Status.State != v1alpha1.BatchFailed &&
		job.Status.State != v1alpha1.BatchRunning &&
		job.Status.State != v1alpha1.BatchSucceeded {
		return fmt.Errorf("InferenceJob %s has not started yet", id)
	}

	// Collect logs from all containers in execution order
	for _, containerName := range jobContainerNames {
		err := utils.StreamLogs(
			c.cfg,
			c.namespace,
			job.Status.PodName,
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