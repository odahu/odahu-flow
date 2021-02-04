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

package packaging

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	odahu_errs "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"

	odahuflowv1alpha1 "github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/packaging"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	mp_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/packaging"
	"github.com/odahu/odahu-flow/packages/operator/pkg/repository/util/kubernetes"
	"github.com/xeipuuv/gojsonschema"
	"go.uber.org/multierr"
)

const (
	ValidationMpErrorMessage             = "Validation of model packaging is failed"
	TrainingIDOrArtifactNameErrorMessage = "empty artifactName"
	ArgumentValidationErrorMessage       = "argument validation is failed: %s"
	EmptyIntegrationNameErrorMessage     = "empty integrationName"
	TargetNotFoundErrorMessage           = "cannot find %s target in packaging integration %s"
	NotValidConnTypeErrorMessage         = "%s target has not valid connection type %s for packaging integration %s"
	UnknownNodeSelector                  = "node selector %v is not presented in ODAHU config"
)

type MpValidator struct {
	piRepo          mp_repository.PackagingIntegrationRepository
	connRepo        conn_repository.Repository
	gpuResourceName string
	packagingConfig config.ModelPackagingConfig
}

func NewMpValidator(
	piRepo mp_repository.PackagingIntegrationRepository,
	connRepo conn_repository.Repository,
	packagingConfig config.ModelPackagingConfig,
	gpuResourceName string,
) *MpValidator {
	return &MpValidator{
		piRepo:          piRepo,
		connRepo:        connRepo,
		packagingConfig: packagingConfig,
		gpuResourceName: gpuResourceName,
	}
}

// ValidateAndSetDefaultsPIRequired validates ModelPackaging fields using a corresponding PackagingIntegration
func (mpv *MpValidator) ValidateAndSetDefaultsPIRequired(mp *packaging.ModelPackaging) (err error) {

	pi, piErr := mpv.piRepo.GetPackagingIntegration(mp.Spec.IntegrationName)
	switch {
	case piErr == nil:
		if len(mp.Spec.Image) == 0 {
			mp.Spec.Image = pi.Spec.DefaultImage
			logMP.Info("Model packaging .spec.image is empty. Set a packaging integration .spec.defaultImage",
				"id", mp.ID, "image", pi.Spec.DefaultImage)
		}
		err = multierr.Append(err, mpv.validateArguments(pi, mp))
		err = multierr.Append(err, mpv.validateTargets(pi, mp))
	case odahu_errs.IsNotFoundError(piErr):
		err = multierr.Append(err,
			fmt.Errorf(
				"packaging integration with name .spec.integrationName = \"%s\" is not found",
				mp.Spec.IntegrationName,
			))
	default:
		// We are not able to complete validation in this case
		panic(fmt.Errorf(
			"unable to get Packaging Integration: %v", piErr,
		))
	}

	return err
}

func (mpv *MpValidator) ValidateAndSetDefaults(mp *packaging.ModelPackaging) (err error) {

	err = multierr.Append(err, mpv.validateMainParameters(mp))
	err = multierr.Append(err, mpv.validateOutputConnection(mp))
	err = multierr.Append(err, mpv.validateNodeSelector(mp))
	err = multierr.Append(err, validation.ValidateResources(mp.Spec.Resources, config.NvidiaResourceName))

	if len(mp.Spec.IntegrationName) == 0 {
		err = multierr.Append(err, errors.New(EmptyIntegrationNameErrorMessage))
	} else {
		err = multierr.Append(err, mpv.ValidateAndSetDefaultsPIRequired(mp))
	}
	if err != nil {
		return fmt.Errorf("%s: %s", ValidationMpErrorMessage, err.Error())
	}
	return nil
}

func (mpv *MpValidator) validateMainParameters(mp *packaging.ModelPackaging) (err error) {
	err = multierr.Append(err, validation.ValidateID(mp.ID))

	if len(mp.Spec.ArtifactName) == 0 {
		err = multierr.Append(err, errors.New(TrainingIDOrArtifactNameErrorMessage))
	}

	if mp.Spec.Resources == nil {
		logMP.Info("Packaging resources parameter is nil. Set the default value",
			"id", mp.ID, "resources", mpv.packagingConfig.DefaultResources)
		mp.Spec.Resources = mpv.packagingConfig.DefaultResources.DeepCopy()
	} else {
		_, resValidationErr := kubernetes.ConvertOdahuflowResourcesToK8s(mp.Spec.Resources, mpv.gpuResourceName)
		err = multierr.Append(err, resValidationErr)
	}

	return err
}

func (mpv *MpValidator) validateArguments(pi *packaging.PackagingIntegration, mp *packaging.ModelPackaging) error {
	required := make([]string, 0)
	if pi.Spec.Schema.Arguments.Required != nil {
		required = pi.Spec.Schema.Arguments.Required
	}

	properties := make(map[string]map[string]interface{})
	if pi.Spec.Schema.Arguments.Properties != nil {
		for _, prop := range pi.Spec.Schema.Arguments.Properties {
			params := make(map[string]interface{})
			for _, param := range prop.Parameters {
				params[param.Name] = param.Value
			}

			properties[prop.Name] = params
		}
	}

	jsonSchema := map[string]interface{}{
		"type":                 "object",
		"properties":           properties,
		"required":             required,
		"additionalProperties": false,
	}

	schemaLoader := gojsonschema.NewGoLoader(jsonSchema)
	data := make(map[string]interface{})
	if mp.Spec.Arguments != nil {
		data = mp.Spec.Arguments
	}
	dataLoader := gojsonschema.NewGoLoader(data)
	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		return err
	}

	if result.Valid() {
		return nil
	}

	return fmt.Errorf(ArgumentValidationErrorMessage, result.Errors())
}

func (mpv *MpValidator) validateTargets(pi *packaging.PackagingIntegration, mp *packaging.ModelPackaging) (err error) {
	requiredTargets := make(map[string]odahuflowv1alpha1.TargetSchema)
	allTargets := make(map[string]odahuflowv1alpha1.TargetSchema)

	for _, target := range pi.Spec.Schema.Targets {
		allTargets[target.Name] = target

		if target.Required {
			requiredTargets[target.Name] = target
		}
	}

	for _, target := range mp.Spec.Targets {
		delete(requiredTargets, target.Name)

		if targetSchema, ok := allTargets[target.Name]; !ok {
			err = multierr.Append(err, fmt.Errorf(TargetNotFoundErrorMessage, target.Name, pi.ID))
		} else {
			delete(allTargets, target.Name)
			if conn, k8sErr := mpv.connRepo.GetConnection(target.ConnectionName); k8sErr != nil {
				err = multierr.Append(err, k8sErr)
			} else {
				isValidConnectionType := false

				for _, connType := range targetSchema.ConnectionTypes {
					if odahuflowv1alpha1.ConnectionType(connType) == conn.Spec.Type {
						isValidConnectionType = true
						break
					}
				}

				if !isValidConnectionType {
					err = multierr.Append(err, fmt.Errorf(NotValidConnTypeErrorMessage, target.Name, conn.Spec.Type, pi.ID))
				}
			}
		}
	}

	// Propagate default values for the remaining targets
	for _, target := range allTargets {
		if len(target.Default) != 0 {
			delete(requiredTargets, target.Name)

			mp.Spec.Targets = append(mp.Spec.Targets, odahuflowv1alpha1.Target{
				Name:           target.Name,
				ConnectionName: target.Default,
			})
		}
	}

	if len(requiredTargets) != 0 {
		requiredTargetsList := make([]string, 0)
		for targetName := range requiredTargets {
			requiredTargetsList = append(requiredTargetsList, targetName)
		}

		err = multierr.Append(err, fmt.Errorf("%s are required targets", requiredTargetsList))
	}

	return err
}

func (mpv *MpValidator) validateOutputConnection(mp *packaging.ModelPackaging) (err error) {
	if len(mp.Spec.OutputConnection) == 0 {
		if len(mpv.packagingConfig.OutputConnectionID) > 0 {
			mp.Spec.OutputConnection = mpv.packagingConfig.OutputConnectionID
			logMP.Info("OutputConnection is empty. Use connection from configuration")
		} else {
			logMP.Info("OutputConnection is empty. Configuration doesn't contain default value")
		}
	}

	emptyErr := validation.ValidateEmpty("OutputConnection", mp.Spec.OutputConnection)
	if emptyErr != nil {
		err = multierr.Append(err, emptyErr)
	}

	notExistsErr := validation.ValidateExistsInRepository(mp.Spec.OutputConnection, mpv.connRepo)
	if notExistsErr != nil {
		err = multierr.Append(err, notExistsErr)
	}

	if err != nil {
		return fmt.Errorf(validation.SpecSectionValidationFailedMessage, "OutputConnection", err.Error())
	}

	return
}

func (mpv *MpValidator) validateNodeSelector(mp *packaging.ModelPackaging) error {
	if len(mp.Spec.NodeSelector) == 0 {
		return nil
	}

	nodePools := mpv.packagingConfig.NodePools

	for _, nodePool := range nodePools {
		if reflect.DeepEqual(mp.Spec.NodeSelector, nodePool.NodeSelector) {
			return nil
		}
	}

	return fmt.Errorf(UnknownNodeSelector, mp.Spec.NodeSelector)
}
