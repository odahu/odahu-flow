/*
 * Copyright 2019 EPAM Systems
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

package servicecatalog

import (
	"bytes"
	"encoding/json"
	model_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"go.uber.org/zap"
	"sync"
	"text/template"
)

type ModelRouteCatalog struct {
	sync.RWMutex
	routeMap map[string]Route
	// Key Deployment ID
	modelsMap map[string]model_types.DeployedModel
	log *zap.SugaredLogger
}

func NewModelRouteCatalog(log *zap.SugaredLogger) *ModelRouteCatalog {
	return &ModelRouteCatalog{
		routeMap: map[string]Route{},
		modelsMap: map[string]model_types.DeployedModel{},
		log: log,
	}
}

// PrefixSwaggerUrls prefix all URLs in swagger using prefix
func PrefixSwaggerUrls(prefix string, swagger model_types.Swagger2) (result model_types.Swagger2, err error) {

	swaggerMap := map[string]interface{}{}
	if err = json.Unmarshal(swagger.Raw, &swaggerMap); err != nil {
		return result, err
	}

	paths := swaggerMap["paths"].(map[string]interface{})

	prefixedPaths := make(map[string]interface{})

	for origURL, data := range paths {
		prefixedURL := prefix +origURL
		prefixedPaths[prefixedURL] = data
	}

	swaggerMap["paths"] = prefixedPaths

	result.Raw, err = json.Marshal(swaggerMap)

	return result, err

}

// TagSwaggerMethods add tags into swagger
func TagSwaggerMethods(tags []string, swagger model_types.Swagger2) (result model_types.Swagger2, err error) {

	swaggerRaw := swagger.Raw

	swaggerMap := map[string]interface{}{}
	if err = json.Unmarshal(swaggerRaw, &swaggerMap); err != nil {
		return result, err
	}
	paths := swaggerMap["paths"].(map[string]interface{})

	for _, method := range paths {

		realMethod := method.(map[string]interface{})
		for _, content := range realMethod {
			realContent := content.(map[string]interface{})
			realContent["tags"] = tags
		}
	}

	result.Raw, err = json.Marshal(swaggerMap)

	return result, err

}

// CreateOrUpdate create or update route in catalog
// All URLs of original swagger of model behind the route will be prefixed by route.Prefix
func (mdc *ModelRouteCatalog) CreateOrUpdate(route Route) error {

	mdc.Lock()
	defer mdc.Unlock()

	prefixedSwagger, err := PrefixSwaggerUrls(route.Prefix, route.Model.ServedModel.Swagger)
	if err != nil {
		return err
	}
	route.Model.ServedModel.Swagger = prefixedSwagger
	mdc.routeMap[route.ID] = route
	mdc.modelsMap[route.Model.DeploymentID] = route.Model
	return nil
}

// Delete delete route from catalog
func (mdc *ModelRouteCatalog) Delete(routeID string, log *zap.SugaredLogger) {
	mdc.Lock()
	defer mdc.Unlock()
	_, ok := mdc.routeMap[routeID]
	if ok {
		delete(mdc.routeMap, routeID)
		log.Info("Model route was deleted")
	}
}

// GetDeployedModel returns information about deployed model
func (mdc *ModelRouteCatalog) GetDeployedModel(deploymentID string) (model_types.DeployedModel, error) {
	mdc.RLock()
	defer mdc.RUnlock()
	model, ok := mdc.modelsMap[deploymentID]
	if !ok {
		return model, odahu_errors.NotFoundError{Entity: deploymentID}
	}
	return model, nil
}

// ProcessSwaggerJSON combine URLs of all models in catalog. It separates endpoints by tagging them using
// route.ID
func (mdc *ModelRouteCatalog) ProcessSwaggerJSON() (string, error) {
	mdc.RLock()
	defer mdc.RUnlock()
	allURLs := map[string]interface{}{}

	for _, route := range mdc.routeMap {

		logger := mdc.log.With("route.id", route.ID)

		taggedSwagger, err := TagSwaggerMethods([]string{route.ID}, route.Model.ServedModel.Swagger)
		if err != nil {
			logger.Errorw("Unable to tag route swagger urls", err)
			continue
		}
		swaggerMap := map[string]interface{}{}
		if err := json.Unmarshal(taggedSwagger.Raw, &swaggerMap); err != nil {
			logger.Errorw("Unable to unmarshall route swagger", err)
			continue
		}

		paths := swaggerMap["paths"].(map[string]interface{})

		for url, data := range paths {
			allURLs[url] = data
		}
	}

	routesBytes, err := json.Marshal(allURLs)
	if err != nil {
		return "", err
	}

	var buff bytes.Buffer
	err = swaggerTemplate.Execute(&buff, string(routesBytes))
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

func init() {
	tmpl, err := template.New("swagger template").Parse(templateStr)
	if err != nil {
		panic(err)
	}

	swaggerTemplate = tmpl
}

const templateStr = `
{
    "swagger": "2.0",
    "info": {
        "description": "Catalog of model services",
        "title": "Service Catalog",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "1.0"
    },
    "schemes": [
      "https"
    ],
    "host": "",
    "basePath": "",
    "paths": {{ . }}
}
`

var swaggerTemplate *template.Template
