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

package modeltraining

import (
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"path"

	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/pkg/apis/odahuflow/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	operator_conf "github.com/odahu/odahu-flow/packages/operator/pkg/config/operator"
	training_conf "github.com/odahu/odahu-flow/packages/operator/pkg/config/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/spf13/viper"
	tektonv1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

const (
	pathToTrainerBin  = "/opt/odahu-flow/trainer"
	modelValidatorBin = "odahuflowctl"
	workspacePath     = "/workspace"
	outputDir         = "output"
	configVolumeName  = "config"
	configDir         = "/etc/odahuflow/"
	configFileName    = "config.yaml"
	configSecretName  = "odahu-flow-training-config" //nolint:gosec
)

func generateTrainerTaskSpec(
	trainingCR *odahuflowv1alpha1.ModelTraining,
	toolchainIntegration *training.ToolchainIntegration,
) (*tektonv1alpha1.TaskSpec, error) {
	mtResources, err := kubernetes.ConvertOdahuflowResourcesToK8s(trainingCR.Spec.Resources)
	if err != nil {
		log.Error(err, "The training resources is not valid",
			"mt_id", trainingCR.Name, "resources", mtResources)

		return nil, err
	}

	helperContainerResources := utils.CalculateHelperContainerResources(mtResources)
	return &tektonv1alpha1.TaskSpec{
		Steps: []tektonv1alpha1.Step{
			createInitTrainerStep(helperContainerResources, trainingCR.Name),
			createMainTrainerStep(trainingCR, toolchainIntegration, &mtResources),
			createArtifactValidationStep(helperContainerResources, trainingCR),
			createResultTrainerStep(helperContainerResources, trainingCR),
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

func createInitTrainerStep(res corev1.ResourceRequirements, mtID string) tektonv1alpha1.Step {
	return tektonv1alpha1.Step{
		Container: corev1.Container{
			Name:    odahuflow.TrainerSetupStep,
			Image:   viper.GetString(training_conf.ModelBuilderImage),
			Command: []string{pathToTrainerBin},
			Args: []string{
				"setup",
				"--mt-file",
				path.Join(workspacePath, mtConfig),
				"--mt-id",
				mtID,
				"--output-connection-name",
				viper.GetString(training_conf.OutputConnectionName),
				"--api-url",
				viper.GetString(operator_conf.APIURL),
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

func createMainTrainerStep(
	train *odahuflowv1alpha1.ModelTraining,
	trainingIntegration *training.ToolchainIntegration,
	trainResources *corev1.ResourceRequirements) tektonv1alpha1.Step {

	envs := make([]corev1.EnvVar, 0, len(trainingIntegration.Spec.AdditionalEnvironments))
	for name, value := range trainingIntegration.Spec.AdditionalEnvironments {
		envs = append(envs, corev1.EnvVar{
			Name:  name,
			Value: value,
		})
	}

	return tektonv1alpha1.Step{
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

func createArtifactValidationStep(
	validatorResources corev1.ResourceRequirements, trainingCR *odahuflowv1alpha1.ModelTraining,
) tektonv1alpha1.Step {
	return tektonv1alpha1.Step{
		Container: corev1.Container{
			Name:    odahuflow.TrainerValidationStep,
			Image:   trainingCR.Spec.Image,
			Command: []string{modelValidatorBin},
			Args: []string{
				"gppi",
				"-m",
				path.Join(workspacePath, outputDir),
				"test",
			},
			Resources: validatorResources,
		},
	}
}

func createResultTrainerStep(
	res corev1.ResourceRequirements, mt *odahuflowv1alpha1.ModelTraining,
) tektonv1alpha1.Step {
	return tektonv1alpha1.Step{
		Container: corev1.Container{
			Name:    odahuflow.TrainerResultStep,
			Image:   viper.GetString(training_conf.ModelBuilderImage),
			Command: []string{pathToTrainerBin},
			Args: []string{
				"result",
				"--mt-file",
				path.Join(workspacePath, mtConfig),
				"--mt-id",
				mt.Name,
				"--output-connection-name",
				mt.Spec.OutputConnection,
				"--api-url",
				viper.GetString(operator_conf.APIURL),
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
