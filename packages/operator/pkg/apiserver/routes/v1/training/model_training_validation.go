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

package training

import (
	"errors"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	mt_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	"go.uber.org/multierr"
	"reflect"
)

const (
	MtVcsNotExistsErrorMessage    = "cannot find VCS Connection"
	EmptyModelNameErrorMessage    = "model name must be non-empty"
	EmptyModelVersionErrorMessage = "model version must be non-empty"
	EmptyVcsNameMessageError      = "VCS name is empty"
	ValidationMtErrorMessage      = "Validation of model training is failed"
	WrongVcsTypeErrorMessage      = "VCS connection must have the GIT type. You pass the connection of %s type"
	WrongVcsReferenceErrorMessage = "you should specify a VCS reference for model training explicitly." +
		" Because %s does not have default reference"
	EmptyDataBindingNameErrorMessage = "you should specify connection name for %d number of data binding"
	EmptyDataBindingPathErrorMessage = "you should specify local path for %d number of data binding"
	WrongDataBindingTypeErrorMessage = "%s data binding has wrong data type. Currently supported the following types" +
		" of connections for data bindings: %v"
	ToolchainEmptyErrorMessage = "toolchain parameter is empty"
	UnknownNodeSelector        = "node selector %v is not presented in ODAHU config"
	defaultIDTemplate          = "%s-%s-%s"
)

var (
	DefaultArtifactOutputTemplate = "{{ .Name }}-{{ .Version }}-{{ .RandomUUID }}.zip"
	expectedConnectionDataTypes   = map[odahuflowv1alpha1.ConnectionType]bool{
		connection.GcsType:       true,
		connection.S3Type:        true,
		connection.AzureBlobType: true,
	}
)

type MtValidator struct {
	mtRepository    mt_repository.ToolchainRepository
	connRepository  conn_repository.Repository
	gpuResourceName string
	trainingConfig  config.ModelTrainingConfig
}

func NewMtValidator(
	mtRepository mt_repository.ToolchainRepository,
	connRepository conn_repository.Repository,
	trainingConfig config.ModelTrainingConfig,
	gpuResourceName string,
) *MtValidator {
	return &MtValidator{
		mtRepository:    mtRepository,
		connRepository:  connRepository,
		trainingConfig:  trainingConfig,
		gpuResourceName: gpuResourceName,
	}
}

func (mtv *MtValidator) ValidatesAndSetDefaults(mt *training.ModelTraining) (err error) {
	err = multierr.Append(err, mtv.validateMainParams(mt))

	err = multierr.Append(err, mtv.validateVCS(mt))

	err = multierr.Append(err, mtv.validateMtData(mt))

	err = multierr.Append(err, mtv.validateToolchain(mt))

	err = multierr.Append(err, mtv.validateOutputConnection(mt))

	err = multierr.Append(err, mtv.validateNodeSelector(mt))

	err = multierr.Append(err, validation.ValidateResources(mt.Spec.Resources, config.NvidiaResourceName))

	if err != nil {
		return fmt.Errorf("%s: %s", ValidationMtErrorMessage, err.Error())
	}

	return
}

func (mtv *MtValidator) validateTrainingVersion(version string) (err error) {
	if len(version) == 0 {
		err = multierr.Append(err, errors.New(EmptyModelVersionErrorMessage))
	}
	if labelErr := validation.ValidateK8sLabel(version); labelErr != nil {
		err = multierr.Append(err, errors.New("invalid training version: "+labelErr.Error()))
	}
	return err
}

func (mtv *MtValidator) validateTrainingName(name string) (err error) {
	if len(name) == 0 {
		err = multierr.Append(err, errors.New(EmptyModelNameErrorMessage))
	}
	if labelErr := validation.ValidateK8sLabel(name); labelErr != nil {
		err = multierr.Append(err, errors.New("invalid training name: "+labelErr.Error()))
	}
	return err
}

func (mtv *MtValidator) validateMainParams(mt *training.ModelTraining) (err error) {
	err = multierr.Append(err, mtv.validateTrainingName(mt.Spec.Model.Name))
	err = multierr.Append(err, mtv.validateTrainingVersion(mt.Spec.Model.Version))

	if len(mt.ID) == 0 {
		u4, uuidErr := uuid.NewV4()
		if uuidErr != nil {
			err = multierr.Append(err, uuidErr)
		} else {
			mt.ID = fmt.Sprintf(defaultIDTemplate, mt.Spec.Model.Name, mt.Spec.Model.Version, u4.String())
			logMT.Info("Training id is empty. Generate a default value", "id", mt.ID)
		}
	}

	err = multierr.Append(err, validation.ValidateID(mt.ID))

	if len(mt.Spec.Model.ArtifactNameTemplate) == 0 {
		logMT.Info("Artifact output template is empty. Set the default value",
			"name", mt.ID, "artifact ame", DefaultArtifactOutputTemplate)
		mt.Spec.Model.ArtifactNameTemplate = DefaultArtifactOutputTemplate
	}

	if mt.Spec.Resources == nil {
		logMT.Info("Training resource parameter is nil. Set the default value",
			"name", mt.ID, "resources", mtv.trainingConfig.DefaultResources)
		mt.Spec.Resources = mtv.trainingConfig.DefaultResources.DeepCopy()
	} else {
		_, resValidationErr := kubernetes.ConvertOdahuflowResourcesToK8s(mt.Spec.Resources, mtv.gpuResourceName)
		err = multierr.Append(err, resValidationErr)
	}

	return err
}

func (mtv *MtValidator) validateToolchain(mt *training.ModelTraining) (err error) {
	if len(mt.Spec.Toolchain) == 0 {
		err = multierr.Append(err, errors.New(ToolchainEmptyErrorMessage))

		return
	}

	if labelErr := validation.ValidateK8sLabel(mt.Spec.Toolchain); labelErr != nil {
		err = multierr.Append(err, errors.New("invalid training toolchain: "+labelErr.Error()))
	}

	toolchain, k8sErr := mtv.mtRepository.GetToolchainIntegration(mt.Spec.Toolchain)
	if k8sErr != nil {
		err = multierr.Append(err, k8sErr)

		return
	}

	if len(mt.Spec.Image) == 0 {
		logMT.Info("Toolchain image parameter is nil. Set the default value",
			"name", mt.ID, "image", toolchain.Spec.DefaultImage)
		mt.Spec.Image = toolchain.Spec.DefaultImage
	}

	return
}

func (mtv *MtValidator) validateVCS(mt *training.ModelTraining) (err error) {
	if len(mt.Spec.VCSName) == 0 {
		err = multierr.Append(err, errors.New(EmptyVcsNameMessageError))

		return
	}

	if vcs, odahuErr := mtv.connRepository.GetConnection(mt.Spec.VCSName); odahuErr != nil {
		logMT.Error(err, MtVcsNotExistsErrorMessage)

		err = multierr.Append(err, odahuErr)
	} else if len(mt.Spec.Reference) == 0 {
		switch {
		case vcs.Spec.Type != connection.GITType:
			err = multierr.Append(err, fmt.Errorf(WrongVcsTypeErrorMessage, vcs.Spec.Type))
		case len(vcs.Spec.Reference) == 0:
			err = multierr.Append(err, fmt.Errorf(WrongVcsReferenceErrorMessage, vcs.ID))
		default:
			logMT.Info("VCS reference parameter is nil. Set the default value",
				"name", mt.ID, "reference", vcs.Spec.Reference)
			mt.Spec.Reference = vcs.Spec.Reference
		}
	}

	return
}

func (mtv *MtValidator) validateMtData(mt *training.ModelTraining) (err error) {
	for i, dbd := range mt.Spec.Data {
		if len(dbd.LocalPath) == 0 {
			err = multierr.Append(err, fmt.Errorf(EmptyDataBindingPathErrorMessage, i))
		}

		if len(dbd.Connection) == 0 {
			err = multierr.Append(err, fmt.Errorf(EmptyDataBindingNameErrorMessage, i))

			continue
		}

		conn, k8sErr := mtv.connRepository.GetConnection(dbd.Connection)
		if k8sErr != nil {
			err = multierr.Append(err, k8sErr)

			continue
		}

		if _, ok := expectedConnectionDataTypes[conn.Spec.Type]; !ok {
			err = multierr.Append(err, fmt.Errorf(
				WrongDataBindingTypeErrorMessage,
				conn.ID, reflect.ValueOf(expectedConnectionDataTypes).MapKeys(),
			))
		}

	}

	return
}

func (mtv *MtValidator) validateOutputConnection(mt *training.ModelTraining) (err error) {
	if len(mt.Spec.OutputConnection) == 0 {
		if len(mtv.trainingConfig.OutputConnectionID) > 0 {
			mt.Spec.OutputConnection = mtv.trainingConfig.OutputConnectionID
			logMT.Info("OutputConnection is empty. Use connection from configuration")
		} else {
			logMT.Info("OutputConnection is empty. Configuration doesn't contain default value")
		}
	}

	emptyErr := validation.ValidateEmpty("OutputConnection", mt.Spec.OutputConnection)
	if emptyErr != nil {
		err = multierr.Append(err, emptyErr)
	}

	notExistsErr := validation.ValidateExistsInRepository(mt.Spec.OutputConnection, mtv.connRepository)
	if notExistsErr != nil {
		err = multierr.Append(err, notExistsErr)
	}

	if err != nil {
		return fmt.Errorf(validation.SpecSectionValidationFailedMessage, "OutputConnection", err.Error())
	}

	return

}

func (mtv *MtValidator) validateNodeSelector(mt *training.ModelTraining) error {
	if len(mt.Spec.NodeSelector) == 0 {
		return nil
	}

	selectorFound := false
	var nodePools []config.NodePool
	if mt.Spec.IsGPUResourceSet() {
		nodePools = mtv.trainingConfig.GPUNodePools
	} else {
		nodePools = mtv.trainingConfig.NodePools
	}

NodePoolsLoop:
	for _, nodePool := range nodePools {
		if len(nodePool.NodeSelector) != len(mt.Spec.NodeSelector) {
			continue
		}
		for key, value := range nodePool.NodeSelector {
			if mt.Spec.NodeSelector[key] != value {
				continue NodePoolsLoop
			}
		}
		selectorFound = true
	}

	if !selectorFound {
		return fmt.Errorf(UnknownNodeSelector, mt.Spec.NodeSelector)
	}
	return nil
}
