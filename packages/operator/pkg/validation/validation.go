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

package validation

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	connection "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	"log"
)

const (
	SpecSectionValidationFailedMessage = "\"Spec.%q\" validation errors: %s"
	EmptyValueStringError              = "%q parameter must be not empty"
)

func ValidateEmpty(parameterName, value string) error {
	if len(value) == 0 {
		return fmt.Errorf(EmptyValueStringError, parameterName)
	}
	return nil
}

func ValidateExistsInRepository(name string, repository connection.Repository) error {
	if len(name) > 0 {
		if _, odahuError := repository.GetConnection(name); odahuError != nil {
			return odahuError
		}
	}
	return nil
}

func InitValidator() {

	translator := en.New()
	uni := ut.New(translator, translator)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
			return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		})

	}
}
