package predictv2

import (
	"encoding/json"
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/predict_v2"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)



type Validator func (content []byte) error

func ValidateInput(content []byte) error {
	var request predict_v2.InferenceRequest
	return json.Unmarshal(content, &request)
}
func ValidateOutput(content []byte) error {
	var response predict_v2.InferenceResponse
	return json.Unmarshal(content, &response)
}

// ValidateDir go through `source` directory files (not recursively) and tries to validate each file
// with .json extension using Validator
// If destination is not nil then files that passed validation will be copied to destination
func ValidateDir(source string, validate Validator, destination *string) ([]string, error) {

	// Files that passed validation
	var passed []string

	if destination != nil {
		if _, err := os.Stat(*destination); os.IsNotExist(err) {
			zap.S().Infof("destination directory specified by -d %s does not exist", *destination)
			if err := os.MkdirAll(*destination, os.ModePerm); err != nil {
				return passed, fmt.Errorf("unable to create %s. error %+v", *destination, err)
			}
			zap.S().Infof("destination directory %s is created", *destination)
		}
	}

	items, err := ioutil.ReadDir(source)
	if err != nil {
		return passed, err
	}

	for _, item := range items {
		file := item.Name()
		if item.IsDir() {
			zap.S().Infof("%s is directory (not file). Skip validating", file)
			continue
		}
		if strings.HasSuffix(file, ".json") {
			fp := filepath.Join(source, file)
			data, err := ioutil.ReadFile(fp)
			if err != nil {
				return passed, err
			}
			parseErr := validate(data)
			if parseErr != nil {
				zap.S().Warnw("Unable to Unmarshal file, validation is failed",
					"filepath", fp, "unmarshal error", parseErr)
				continue
			}

			if destination != nil {
				if err := ioutil.WriteFile(filepath.Join(*destination, file), data, item.Mode()); err != nil {
					return passed, err
				}

			}
			zap.S().Infof("%s successfully passed validation", file)
			passed = append(passed, file)


		} else {
			zap.S().Infof("%s has not .json extension. Skip validating", file)
		}
	}

	return passed, nil
}