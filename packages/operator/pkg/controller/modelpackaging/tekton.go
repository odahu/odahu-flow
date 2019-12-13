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

package modelpackaging

import (
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/pkg/apis/odahuflow/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	operator_conf "github.com/odahu/odahu-flow/packages/operator/pkg/config/operator"
	packaging_conf "github.com/odahu/odahu-flow/packages/operator/pkg/config/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/odahuflow"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/spf13/viper"
	tektonv1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"path"
)

const (
	pathToPackagerBin = "/opt/odahu-flow/packager"
	workspacePath     = "/workspace"
	outputDir         = "output"
	configVolumeName  = "config"
	configDir         = "/etc/odahuflow/"
	configFileName    = "config.yaml"
	configSecretName  = "odahu-flow-packaging-config" //nolint:gosec
)

func generatePackagerTaskSpec(
	packagingCR *odahuflowv1alpha1.ModelPackaging,
	packagingIntegration *packaging.PackagingIntegration,
) (*tektonv1alpha1.TaskSpec, error) {
	mainPackagerStep, err := createMainPackagerStep(packagingCR, packagingIntegration)
	if err != nil {
		return nil, err
	}

	return &tektonv1alpha1.TaskSpec{
		Steps: []tektonv1alpha1.Step{
			createInitPackagerStep(packagingCR),
			mainPackagerStep,
			createResultPackagerStep(packagingCR.Name),
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

func createInitPackagerStep(packagingCR *odahuflowv1alpha1.ModelPackaging) tektonv1alpha1.Step {
	return tektonv1alpha1.Step{
		Container: corev1.Container{
			Name:    odahuflow.PackagerSetupStep,
			Image:   viper.GetString(packaging_conf.ModelPackagerImage),
			Command: []string{pathToPackagerBin},
			Args: []string{
				"setup",
				"--mp-file",
				path.Join(workspacePath, mpContentFile),
				"--mp-id",
				packagingCR.Name,
				"--output-connection-name",
				packagingCR.Spec.OutputConnection,
				"--api-url",
				viper.GetString(operator_conf.APIURL),
				"--config",
				path.Join(configDir, configFileName),
			},
			Resources: packagerResources,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      configVolumeName,
					MountPath: configDir,
				},
			},
		},
	}
}

func createMainPackagerStep(
	packagingCR *odahuflowv1alpha1.ModelPackaging,
	packagingIntegration *packaging.PackagingIntegration) (tektonv1alpha1.Step, error) {
	packResources, err := kubernetes.ConvertOdahuflowResourcesToK8s(packagingCR.Spec.Resources)
	if err != nil {
		log.Error(err, "The packaging resources is not valid",
			"mp id", packagingCR.Name, "resources", packagingCR.Namespace)

		return tektonv1alpha1.Step{}, err
	}

	return tektonv1alpha1.Step{
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

func createResultPackagerStep(mpID string) tektonv1alpha1.Step {
	return tektonv1alpha1.Step{
		Container: corev1.Container{
			Name:    odahuflow.PackagerResultStep,
			Image:   viper.GetString(packaging_conf.ModelPackagerImage),
			Command: []string{pathToPackagerBin},
			Args: []string{
				"result",
				"--mp-file",
				path.Join(workspacePath, mpContentFile),
				"--mp-id",
				mpID,
				"--output-connection-name",
				viper.GetString(packaging_conf.OutputConnectionName),
				"--api-url",
				viper.GetString(operator_conf.APIURL),
				"--config",
				path.Join(configDir, configFileName),
			},
			Resources: packagerResources,
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      configVolumeName,
					MountPath: configDir,
				},
			},
		},
	}
}
