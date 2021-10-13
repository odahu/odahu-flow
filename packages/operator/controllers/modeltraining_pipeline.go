/*
 * Copyright 2019 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controllers

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"path"

	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

const (
	pathToTrainerBin = "/opt/odahu-flow/trainer"
	workspacePath    = "/workspace"
	outputDir        = "output"
	configVolumeName = "config"
	configDir        = "/etc/odahuflow/"
	configFileName   = "config.yaml"
	configSecretName = "odahu-flow-training-config" //nolint:gosec
)

func (r *ModelTrainingReconciler) generateTrainerTaskSpec(
	trainingCR *odahuflowv1alpha1.ModelTraining,
	trainingIntegration *training.TrainingIntegration,
) (*tektonv1beta1.TaskSpec, error) {
	mtResources, err := kubernetes.ConvertOdahuflowResourcesToK8s(trainingCR.Spec.Resources, r.gpuResourceName)
	if err != nil {
		log.Error(err, "The training resources is not valid",
			"mt_id", trainingCR.Name, "resources", mtResources)

		return nil, err
	}

	helperContainerResources := utils.CalculateHelperContainerResources(mtResources, r.gpuResourceName)
	return &tektonv1beta1.TaskSpec{
		Steps: []tektonv1beta1.Step{
			r.createInitTrainerStep(helperContainerResources, trainingCR.Name),
			r.createMainTrainerStep(trainingCR, trainingIntegration, &mtResources),
			r.createResultTrainerStep(helperContainerResources, trainingCR),
		},
		Volumes: []corev1.Volume{
			{
				Name: configVolumeName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: configSecretName,
					},
				},
			},
		},
	}, nil
}

func (r *ModelTrainingReconciler) createInitTrainerStep(
	res corev1.ResourceRequirements, mtID string,
) tektonv1beta1.Step {
	return tektonv1beta1.Step{
		Container: corev1.Container{
			Name:    odahuflow.TrainerSetupStep,
			Image:   r.trainingConfig.ModelTrainerImage,
			Command: []string{pathToTrainerBin},
			Args: []string{
				"setup",
				"--mt-file",
				path.Join(workspacePath, mtConfig),
				"--mt-id",
				mtID,
				"--api-url",
				r.operatorConfig.Auth.APIURL,
				"--config",
				path.Join(configDir, configFileName),
			},
			Resources: res,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      configVolumeName,
					MountPath: configDir,
				},
			},
		},
	}
}

func (r *ModelTrainingReconciler) createMainTrainerStep(
	train *odahuflowv1alpha1.ModelTraining,
	trainingIntegration *training.TrainingIntegration,
	trainResources *corev1.ResourceRequirements) tektonv1beta1.Step {

	envs := make([]corev1.EnvVar, 0, len(trainingIntegration.Spec.AdditionalEnvironments)+len(train.Spec.CustomEnvs))
	for name, value := range trainingIntegration.Spec.AdditionalEnvironments {
		envs = append(envs, corev1.EnvVar{
			Name:  name,
			Value: value,
		})
	}
	for _, env := range train.Spec.CustomEnvs {
		envs = append(envs, corev1.EnvVar{
			Name:  env.Name,
			Value: env.Value,
		})
	}

	return tektonv1beta1.Step{
		Container: corev1.Container{
			Name:    odahuflow.TrainerTrainStep,
			Image:   train.Spec.Image,
			Command: []string{trainingIntegration.Spec.Entrypoint},
			Env:     envs,
			Args: []string{
				"--verbose",
				"--mt",
				path.Join(workspacePath, mtConfig),
				"--target",
				path.Join(workspacePath, outputDir),
			},
			Resources: *trainResources,
		},
	}
}

func (r *ModelTrainingReconciler) createResultTrainerStep(
	res corev1.ResourceRequirements, mt *odahuflowv1alpha1.ModelTraining,
) tektonv1beta1.Step {
	return tektonv1beta1.Step{
		Container: corev1.Container{
			Name:    odahuflow.TrainerResultStep,
			Image:   r.trainingConfig.ModelTrainerImage,
			Command: []string{pathToTrainerBin},
			Args: []string{
				"result",
				"--mt-file",
				path.Join(workspacePath, mtConfig),
				"--mt-id",
				mt.Name,
				"--api-url",
				r.operatorConfig.Auth.APIURL,
				"--output-dir",
				path.Join(workspacePath, outputDir),
				"--config",
				path.Join(configDir, configFileName),
			},
			Resources: res,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      configVolumeName,
					MountPath: configDir,
				},
			},
		},
	}
}
