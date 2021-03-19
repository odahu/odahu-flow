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

package batchinferenceutils

import (
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/rclone"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	. "github.com/odahu/odahu-flow/packages/operator/controllers/types"
	"strings"
)


func GetBucketNamePath(connName string, path string, connAPI ConnGetter) (
	bucketName string, actualPath string, err error) {

	var conn *connection.Connection
	var connectionPath string

	conn, err = connAPI.GetConnection(connName)
	if err != nil {
		return
	}

	bucketName, connectionPath, err = rclone.GetBucketAndPath(&conn.Spec)
	if err != nil {
		return
	}

	if path == "" {
		actualPath = connectionPath
	} else {
		actualPath = path
	}

	return

}

// GetModelSyncSteps return steps that a specific for certain model synchronization case
// depending on v1alpha1.ModelSource (remote or local) and connection type
func GetModelSyncSteps(
	job *v1alpha1.BatchInferenceJob,
	rcloneImage string,
	res corev1.ResourceRequirements,
	connAPI ConnGetter,
) (steps []tektonv1beta1.Step, err error) {
	switch {
	case job.Spec.ModelSource.Remote != nil:
		source := job.Spec.ModelSource.Remote
		bucket, path, err := GetBucketNamePath(source.ModelConnection, source.ModelPath, connAPI)
		if err != nil {
			return steps, err
		}
		if strings.HasSuffix(path, ".tar.gz") || strings.HasSuffix(path, ".zip"){
			steps = append(steps, GetSyncPackedModelStep(
				rcloneImage, source.ModelConnection, bucket, path, res))
		} else{
			steps = append(steps, GetSyncModelStep(rcloneImage, source.ModelConnection, bucket, path, res))
		}
	case job.Spec.ModelSource.Local != nil:}  // Local model source does not require to sync model
	return steps, nil
}

// BatchJobToTaskSpec generate tektoncd TaskSpec based on v1alpha1.BatchInferenceJob
func BatchJobToTaskSpec(job *v1alpha1.BatchInferenceJob,
	connAPI ConnGetter,
	gpuResourceName string, rcloneImage string,
	toolsSecret string,
	toolsImage string,
) (ts *tektonv1beta1.TaskSpec, err error) {

	jobRes, err := kubernetes.ConvertOdahuflowResourcesToK8s(job.Spec.Resources, gpuResourceName)
	if err != nil {
		return ts, fmt.Errorf("unable to convert odahu resources to kubernetes resources: %s", err)
	}

	helpContainerRes := utils.CalculateHelperContainerResources(jobRes, gpuResourceName)

	bucket, path, err := GetBucketNamePath(job.Spec.InputConnection, job.Spec.InputPath, connAPI)
	if err != nil {
		return ts, err
	}

	connsForRClone := []string{job.Spec.InputConnection, job.Spec.OutputConnection}
	if job.Spec.ModelSource.Remote != nil {
		connsForRClone = append(connsForRClone, job.Spec.ModelSource.Remote.ModelConnection)
	}

	steps := []tektonv1beta1.Step{
		GetConfigureRCloneStep(
			toolsImage,
			helpContainerRes,
			connsForRClone...),
		GetSyncDataStep(rcloneImage, job.Spec.InputConnection, bucket, path, helpContainerRes),
	}

	modelSyncSteps, err := GetModelSyncSteps(job, rcloneImage, helpContainerRes, connAPI)
	if err != nil {
		return ts, err
	}
	steps = append(steps, modelSyncSteps...)

	bucket, path, err = GetBucketNamePath(job.Spec.OutputConnection, job.Spec.OutputPath, connAPI)
	if err != nil {
		return ts, err
	}

	steps = append(steps, []tektonv1beta1.Step{
		GetValidateInputStep(toolsImage, helpContainerRes),
		GetLogInputStep(toolsImage, job.Spec.BatchRequestID, helpContainerRes),
		GetUserContainer(job.Spec.Image, job.Spec.Command, job.Spec.Args, jobRes),
		GetValidateOutputStep(toolsImage, helpContainerRes),
		GetLogOutputStep(toolsImage, job.Spec.BatchRequestID, helpContainerRes),
		GetSyncOutputStep(rcloneImage, job.Spec.OutputConnection, bucket, path, helpContainerRes),
	}...)

	ts = &tektonv1beta1.TaskSpec{
		Steps:   steps,
		Volumes: []corev1.Volume{
			{
				Name: toolsConfigVolume,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: toolsSecret,
					},
				},
			},
		},
	}

	return ts, err
}