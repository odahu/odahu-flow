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

package deployment

import (
	"errors"
	"fmt"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/deployment"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	"go.uber.org/multierr"
	"reflect"
)

const (
	EmptyImageErrorMessage             = "the image parameter is empty"
	NegativeMinReplicasErrorMessage    = "minimum number of replicas parameter must not be less than 0"
	NegativeMaxReplicasErrorMessage    = "maximum number of replicas parameter must not be less than 1"
	MaxMoreThanMinReplicasErrorMessage = "maximum number of replicas parameter must not be less than minimum number " +
		"of replicas parameter"
	ReadinessProbeErrorMessage = "readinessProbeInitialDelay must be non-negative integer"
	LivenessProbeErrorMessage  = "livenessProbeInitialDelay must be non-negative integer"
	UnknownNodeSelector        = "node selector %v is not presented in ODAHU config"
	DefaultRolePrefix		   = "role-"
)

var (
	MdDefaultMinimumReplicas            = int32(0)
	MdDefaultMaximumReplicas            = int32(1)
	MdDefaultLivenessProbeInitialDelay  = int32(0)
	MdDefaultReadinessProbeInitialDelay = int32(0)
)

type ModelDeploymentValidator struct {
	modelDeploymentConfig config.ModelDeploymentConfig
	gpuResourceName       string
	defaultResources      odahuflowv1alpha1.ResourceRequirements
}

func NewModelDeploymentValidator(
	modelDeploymentConfig config.ModelDeploymentConfig,
	gpuResourceName string,
) *ModelDeploymentValidator {
	return &ModelDeploymentValidator{
		modelDeploymentConfig: modelDeploymentConfig,
		gpuResourceName:       gpuResourceName,
		defaultResources:      modelDeploymentConfig.DefaultResources,
	}
}

func (mdv *ModelDeploymentValidator) ValidatesMDAndSetDefaults(md *deployment.ModelDeployment) (err error) {
	err = multierr.Append(err, validation.ValidateID(md.ID))

	if len(md.Spec.Image) == 0 {
		err = multierr.Append(err, errors.New(EmptyImageErrorMessage))
	}

	if md.Spec.RoleName == nil || len(*md.Spec.RoleName) == 0 {
		defaultRoleName := DefaultRolePrefix + md.ID
		logMD.Info("Role name parameter is nil or empty. Set the model Role as the model ID with a prefix",
			"Deployment name", md.ID, "role name", defaultRoleName)
		md.Spec.RoleName = &defaultRoleName
	} else {
		err = multierr.Append(err, validation.ValidateK8sLabel(*md.Spec.RoleName))
	}

	if md.Spec.MinReplicas == nil {
		logMD.Info("Minimum number of replicas parameter is nil. Set the default value",
			"Deployment name", md.ID, "replicas", MdDefaultMinimumReplicas)
		md.Spec.MinReplicas = &MdDefaultMinimumReplicas
	} else if *md.Spec.MinReplicas < 0 {
		err = multierr.Append(errors.New(NegativeMinReplicasErrorMessage), err)
	}

	if md.Spec.MaxReplicas == nil {
		if *md.Spec.MinReplicas > MdDefaultMaximumReplicas {
			md.Spec.MaxReplicas = md.Spec.MinReplicas
		} else {
			md.Spec.MaxReplicas = &MdDefaultMaximumReplicas
		}

		logMD.Info("Maximum number of replicas parameter is nil. Set the default value",
			"Deployment name", md.ID, "replicas", *md.Spec.MinReplicas)
	} else if *md.Spec.MaxReplicas < 1 {
		err = multierr.Append(errors.New(NegativeMaxReplicasErrorMessage), err)
	}

	if md.Spec.MaxReplicas != nil && md.Spec.MinReplicas != nil && *md.Spec.MinReplicas > *md.Spec.MaxReplicas {
		err = multierr.Append(errors.New(MaxMoreThanMinReplicasErrorMessage), err)
	}

	if md.Spec.Resources == nil {
		logMD.Info("Deployment resources parameter is nil. Set the default value",
			"Deployment name", md.ID, "resources", mdv.defaultResources)
		md.Spec.Resources = mdv.defaultResources.DeepCopy()
	} else {
		_, resValidationErr := kubernetes.ConvertOdahuflowResourcesToK8s(md.Spec.Resources, mdv.gpuResourceName)
		err = multierr.Append(err, resValidationErr)
	}

	if md.Spec.ReadinessProbeInitialDelay == nil {
		logMD.Info("readinessProbeInitialDelay parameter is nil. Set the default value",
			"Deployment name", md.ID, "readinessProbeInitialDelay", MdDefaultReadinessProbeInitialDelay)
		md.Spec.ReadinessProbeInitialDelay = &MdDefaultReadinessProbeInitialDelay
	} else if *md.Spec.ReadinessProbeInitialDelay < 0 {
		err = multierr.Append(errors.New(ReadinessProbeErrorMessage), err)
	}

	if md.Spec.LivenessProbeInitialDelay == nil {
		logMD.Info("livenessProbeInitialDelay is nil. Set the default value",
			"Deployment name", md.ID, "livenessProbeInitialDelay", MdDefaultLivenessProbeInitialDelay)

		md.Spec.LivenessProbeInitialDelay = &MdDefaultLivenessProbeInitialDelay
	} else if *md.Spec.LivenessProbeInitialDelay < 0 {
		err = multierr.Append(errors.New(LivenessProbeErrorMessage), err)
	}

	if md.Spec.ImagePullConnectionID == nil || len(*md.Spec.ImagePullConnectionID) == 0 {
		logMD.Info(
			"imagePullConnID parameter is nil. Set the default value",
			"Deployment name", md.ID,
			"imagePullConnID", mdv.modelDeploymentConfig.DefaultDockerPullConnName,
		)

		md.Spec.ImagePullConnectionID = &mdv.modelDeploymentConfig.DefaultDockerPullConnName
	}

	err = multierr.Append(mdv.validateNodeSelector(md), err)

	err = multierr.Append(err, validation.ValidateResources(md.Spec.Resources, config.NvidiaResourceName))

	return err
}

func (mdv *ModelDeploymentValidator) validateNodeSelector(md *deployment.ModelDeployment) error {
	if len(md.Spec.NodeSelector) == 0 {
		return nil
	}

	nodePools := mdv.modelDeploymentConfig.NodePools

	for _, nodePool := range nodePools {
		if reflect.DeepEqual(md.Spec.NodeSelector, nodePool.NodeSelector) {
			return nil
		}
	}

	return fmt.Errorf(UnknownNodeSelector, md.Spec.NodeSelector)
}
