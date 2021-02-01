/*
 * Copyright 2021 EPAM Systems
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

package inspectors

import (
	"bytes"
	"encoding/json"
	"fmt"
	model_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"
	"github.com/odahu/odahu-flow/packages/operator/pkg/servicecatalog/inspectors/bindata"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"path"
	"text/template"
)

const (
	specTemplateFilename = "triton_oas_template.json"
)

type TritonInspector struct {
	EdgeHost   string
	EdgeURL    url.URL
	HTTPClient httpClient
}

func (t TritonInspector) Inspect(prefix string, log *zap.SugaredLogger) (model_types.ServedModel, error) {
	log.Info("getting a list of served models")
	listModelsRequest := t.generateListModelsRequest(prefix)
	response, err := t.HTTPClient.Do(listModelsRequest)
	if err != nil {
		log.Errorw("failed to fetch model repository", "prefix", prefix)
		return model_types.ServedModel{}, fmt.Errorf("failed to fetch model repository on prefix %s", prefix)
	}

	var models []TritonModelMeta
	decoder := json.NewDecoder(response.Body)
	if decoder.Decode(&models) != nil {
		log.Errorw("failed to unmarshall model repository", "prefix", prefix)
		return model_types.ServedModel{}, fmt.Errorf("failed to unmarshall model repository on prefix %s", prefix)
	}
	log.Info("found models", models)

	var model TritonModelMeta

	switch len(models) {
	case 0:
		log.Errorw("triton server serves 0 models", "prefix", prefix)
		return model_types.ServedModel{}, fmt.Errorf("triton server serves 0 models on prefix %s", prefix)
	case 1:
		model = models[0]
	default:
		log.Warnw("triton server serves more than 1 model", "prefix", prefix)
		model = models[0]
	}

	log.Info("rendering spec template")
	specTemplate := bindata.MustAsset(specTemplateFilename)

	tpl, err := template.New("_").Parse(string(specTemplate))
	if err != nil {
		return model_types.ServedModel{}, err
	}

	b := bytes.NewBuffer([]byte{})
	err = tpl.Execute(b, struct {
		ModelName    string
		ModelVersion string
	}{
		ModelName:    model.Name,
		ModelVersion: model.Version,
	})
	if err != nil {
		return model_types.ServedModel{}, err
	}

	specBytes := b.Bytes()
	return model_types.ServedModel{
		Swagger:  model_types.Swagger2{Raw: specBytes},
		Metadata: model_types.Metadata{},
	}, nil
}

func (t *TritonInspector) generateListModelsRequest(prefix string) *http.Request {
	serverUrl := url.URL{
		Scheme: t.EdgeURL.Scheme,
		Host:   t.EdgeURL.Host,
		Path:   path.Join(t.EdgeURL.Path, prefix, "v2/repository/index"),
	}

	return &http.Request{
		Method: http.MethodPost,
		URL:    &serverUrl,
		Host:   t.EdgeHost,
	}
}

type TritonModelMeta struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
