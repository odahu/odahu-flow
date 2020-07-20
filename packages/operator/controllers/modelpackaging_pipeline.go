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

package controllers

import (
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	tektonv1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"path"
)

const (
	pathToPackagerBin         = "/opt/odahu-flow/packager"
	configPackagingSecretName = "odahu-flow-packaging-config" //nolint:gosec
)

func (r *ModelPackagingReconciler) generatePackagerTaskSpec(
	packagingCR *odahuflowv1alpha1.ModelPackaging,
	packagingIntegration *packaging.PackagingIntegration,
) (*tektonv1beta1.TaskSpec, error) {
	mainPackagerStep, err := r.createMainPackagerStep(packagingCR, packagingIntegration)
	if err != nil {
		return nil, err
	}

	helperContainerResources := utils.CalculateHelperContainerResources(mainPackagerStep.Resources, r.gpuResourceName)
	return &tektonv1beta1.TaskSpec{
		Steps: []tektonv1beta1.Step{
			r.createInitPackagerStep(helperContainerResources, packagingCR),
			mainPackagerStep,
			r.createResultPackagerStep(helperContainerResources, packagingCR.Name),
		},
		Volumes: []corev1.Volume{
			{
				Name: configVolumeName,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: configPackagingSecretName,
					},
				},
			},
		},
	}, nil
}

func (r *ModelPackagingReconciler) createInitPackagerStep(
	res corev1.ResourceRequirements, packagingCR *odahuflowv1alpha1.ModelPackaging,
) tektonv1beta1.Step {
	return tektonv1beta1.Step{
		Container: corev1.Container{
			Name:    odahuflow.PackagerSetupStep,
			Image:   r.packagingConfig.ModelPackagerImage,
			Command: []string{pathToPackagerBin},
			Args: []string{
				"setup",
				"--mp-file",
				path.Join(workspacePath, mpContentFile),
				"--mp-id",
				packagingCR.Name,
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

func (r *ModelPackagingReconciler) createMainPackagerStep(
	packagingCR *odahuflowv1alpha1.ModelPackaging,
	packagingIntegration *packaging.PackagingIntegration) (tektonv1beta1.Step, error) {
	packResources, err := kubernetes.ConvertOdahuflowResourcesToK8s(packagingCR.Spec.Resources, r.gpuResourceName)
	if err != nil {
		log.Error(err, "The packaging resources is not valid",
			"mp id", packagingCR.Name, "resources", packagingCR.Namespace)

		return tektonv1beta1.Step{}, err
	}

	return tektonv1beta1.Step{
		Container: corev1.Container{
			Name:    odahuflow.PackagerPackageStep,
			Image:   packagingCR.Spec.Image,
			Command: []string{packagingIntegration.Spec.Entrypoint},
			Args: []string{
				path.Join(workspacePath, outputDir),
				path.Join(workspacePath, mpContentFile),
			},
			SecurityContext: &corev1.SecurityContext{
				Privileged:               &packagingPrivileged,
				AllowPrivilegeEscalation: &packagingPrivileged,
			},
			Resources: packResources,
		},
	}, nil
}

func (r *ModelPackagingReconciler) createResultPackagerStep(
	res corev1.ResourceRequirements, mpID string,
) tektonv1beta1.Step {
	return tektonv1beta1.Step{
		Container: corev1.Container{
			Name:    odahuflow.PackagerResultStep,
			Image:   r.packagingConfig.ModelPackagerImage,
			Command: []string{pathToPackagerBin},
			Args: []string{
				"result",
				"--mp-file",
				path.Join(workspacePath, mpContentFile),
				"--mp-id",
				mpID,
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
