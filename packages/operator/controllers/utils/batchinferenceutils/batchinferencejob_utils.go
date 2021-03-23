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
	"errors"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	. "github.com/odahu/odahu-flow/packages/operator/controllers/types"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/rclone"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/model_registry/object_storage"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
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

func DiscoverNameVersion(
	connName string, modelPath string,
	connType v1alpha1.ConnectionType, connAPI ConnGetter) (name string, version string, err error) {

	switch {
	case connType == connection.GcsType || connType == connection.S3Type || connType == connection.AzureBlobType:
		mr := object_storage.NewModelRegistry(connAPI)
		return mr.Meta(connName, modelPath)
	default:
		return "", "", fmt.Errorf(
			`connection type "%s" is not supported to discover model meta`, connType)
	}
}

// GetModelSyncSteps return steps that a specific for certain model synchronization case
// depending on v1alpha1.ModelSource (remote or local) and connection type
func GetModelSyncSteps(
	job *v1alpha1.BatchInferenceJob,
	res corev1.ResourceRequirements,
	odahuToolsImage string,
	connType v1alpha1.ConnectionType,
) (steps []tektonv1beta1.Step, err error) {

	if job.Spec.ModelSource.Remote == nil {
		return steps, err
	}
	connName := job.Spec.ModelSource.Remote.ModelConnection
	modelPath := job.Spec.ModelSource.Remote.ModelPath

	// Select Model sync algorithm according to connection type
	// Different connection type usually mean different model registries
	switch  {
	case connType == connection.GcsType || connType == connection.S3Type || connType == connection.AzureBlobType:
		steps = append(steps, GetObjectStorageModelSyncStep(odahuToolsImage, connName, modelPath, res))
	default:
		return steps, fmt.Errorf(`connection type "%s" is not supported to model sync`, connType)
	}
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


	steps := []tektonv1beta1.Step{
		GetConfigureRCloneStep(
			toolsImage,
			helpContainerRes,
			[]string{job.Spec.InputConnection, job.Spec.OutputConnection}...),
		GetSyncDataStep(rcloneImage, job.Spec.InputConnection, bucket, path, helpContainerRes),
	}


	// If modelSource is remote than add step that fetch model from model registry to workspace
	// And we need to discover model name and version from this model registry
	var modelName, modelVersion string
	var modelPathEnv corev1.EnvVar
	switch {
	case job.Spec.ModelSource.Remote != nil:

		connName := job.Spec.ModelSource.Remote.ModelConnection
		remoteModelPath := job.Spec.ModelSource.Remote.ModelPath

		conn, err := connAPI.GetConnection(connName)
		if err != nil {
			return ts, err
		}
		connType := conn.Spec.Type

		// Add model sync step to pipeline
		modelSyncSteps, err := GetModelSyncSteps(job, helpContainerRes, toolsImage, connType)
		if err != nil {
			return ts, err
		}
		steps = append(steps, modelSyncSteps...)

		// Discover Model name and version from registry
		modelName, modelVersion, err = DiscoverNameVersion(connName, remoteModelPath, connType, connAPI)
		if err != nil {
			return ts, err
		}
		// Set modelPathEnv to default location (because model from registry always synced to pre-defined folder)
		modelPathEnv = DefaultOdahuModelPathEnv
	case job.Spec.ModelSource.Local != nil:
		modelName = job.Spec.ModelSource.Local.ModelMeta.Name
		modelVersion = job.Spec.ModelSource.Local.ModelMeta.Version
		modelPathEnv = corev1.EnvVar{
			Name:      odahuModelPathEnvName,
			Value:     job.Spec.ModelSource.Local.ModelPath,
		}
	default:
		return ts, errors.New("whether .Spec.ModelSource.Remote or .Spec.ModelSource.Local should be defined")
	}



	bucket, path, err = GetBucketNamePath(job.Spec.OutputConnection, job.Spec.OutputPath, connAPI)
	if err != nil {
		return ts, err
	}

	steps = append(steps, []tektonv1beta1.Step{
		GetValidateInputStep(toolsImage, helpContainerRes),
		GetLogInputStep(toolsImage, job.Spec.BatchRequestID, helpContainerRes, modelName, modelVersion),
		GetUserContainer(job.Spec.Image, job.Spec.Command, job.Spec.Args, jobRes, modelPathEnv),
		GetValidateOutputStep(toolsImage, helpContainerRes),
		GetLogOutputStep(toolsImage, job.Spec.BatchRequestID, helpContainerRes, modelName, modelVersion),
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