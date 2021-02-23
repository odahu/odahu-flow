package predict_v2

import (
	"encoding/json"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/predict_v2"
	"go.uber.org/zap"
	"io/ioutil"
	"path/filepath"
	"strings"
)


type ValidationHook func(dir string, file string, raw []byte, request *predict_v2.InferenceRequest) error

// ValidateDir go through `source` directory files (not recursively) and tries to unmarshal each file
// with .json extension to predict_v2.InferenceRequest
// Error during unmarshal means validation fail
// If destination is not nil then files that passed validation will be copied to destination
func ValidateDir(source string, destination *string) ([]string, error) {


	// Files that passed validation
	var passed []string

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
			var request predict_v2.InferenceRequest
			parseErr := json.Unmarshal(data, &request)
			if parseErr != nil {
				zap.S().Warnw("Unable to Unmarshal file, validation is failed",
					"filepath", fp, "unmarshal error", parseErr)
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