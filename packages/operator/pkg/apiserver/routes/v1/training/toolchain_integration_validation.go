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
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/validation"
	"go.uber.org/multierr"
)

const (
	ValidationTiErrorMessage      = "Validation of toolchain integration is failed"
	EmptyEntrypointErrorMessage   = "entrypoint must be no empty"
	EmptyDefaultImageErrorMessage = "defaultImage must be no empty"
)

type TiValidator struct {
}

func NewTiValidator() *TiValidator {
	return &TiValidator{}
}

func (tiv *TiValidator) ValidatesAndSetDefaults(ti *training.ToolchainIntegration) (err error) {
	err = multierr.Append(err, validation.ValidateID(ti.ID))

	if len(ti.Spec.Entrypoint) == 0 {
		err = multierr.Append(err, errors.New(EmptyEntrypointErrorMessage))
	}

	if len(ti.Spec.DefaultImage) == 0 {
		err = multierr.Append(err, errors.New(EmptyDefaultImageErrorMessage))
	}

	if err != nil {
		return fmt.Errorf("%s: %s", ValidationTiErrorMessage, err.Error())
	}

	return
}
