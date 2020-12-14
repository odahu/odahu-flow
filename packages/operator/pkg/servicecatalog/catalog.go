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
	"go.uber.org/zap"
	"sync"
	"text/template"

	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("catalog")

type ModelRouteInfo struct {
	Data map[string]map[string]interface{}
}

type ModelRouteCatalog struct {
	sync.RWMutex
	routes map[string]*ModelRouteInfo
	routeMap map[string]route
	log *zap.SugaredLogger
}

func NewModelRouteCatalog() *ModelRouteCatalog {
	return &ModelRouteCatalog{
		routes: map[string]*ModelRouteInfo{},
	}
}

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

func (mdc *ModelRouteCatalog) AddModelRoute(mr *v1alpha1.ModelRoute, infoResponse []byte) error {
	mdc.Lock()
	defer mdc.Unlock()
	log.Info("Add model route", "model route id", mr.Name)

	modelSwagger := map[string]interface{}{}
	if err := json.Unmarshal(infoResponse, &modelSwagger); err != nil {
		log.Error(err, "Unmarshal swagger model", "mr id", mr.Name)
		return err
	}
	paths := modelSwagger["paths"].(map[string]interface{})

	for url, method := range paths {
		realURL := mr.Spec.URLPrefix + url

		realMethod := method.(map[string]interface{})
		for _, content := range realMethod {
			realContent := content.(map[string]interface{})
			realContent["tags"] = []string{mr.Name}
		}

		log.Info("Add model url", "model url", realURL, "model route id", mr.Name)

		mri, ok := mdc.routes[mr.Name]
		if !ok {
			mri = &ModelRouteInfo{}
			mri.Data = make(map[string]map[string]interface{})
			mdc.routes[mr.Name] = mri
		}

		mri.Data[realURL] = realMethod
	}

	return nil
}

// Route's model swagger modified by prefixing using route prefix
// And add to local store (index of routes)
func (mdc *ModelRouteCatalog) CreateOrUpdate(route route) error {

	mdc.Lock()
	defer mdc.Unlock()

	prefixedSwagger, err := PrefixSwaggerUrls(route.prefix, route.model.ServedModel.Swagger)
	if err != nil {
		return err
	}
	route.model.ServedModel.Swagger = prefixedSwagger
	mdc.routeMap[route.id] = route
	return nil
}


func (mdc *ModelRouteCatalog) Delete(routeID string) {
	mdc.Lock()
	defer mdc.Unlock()
	delete(mdc.routeMap, routeID)
}

func (mdc *ModelRouteCatalog) DeleteModelRoute(mrName string) {
	mdc.Lock()
	defer mdc.Unlock()
	log.Info("Delete model route", "model route id", mrName)

	delete(mdc.routes, mrName)
}

// Combine swagger `paths` of different models into single swagger page
// And return its content
func (mdc *ModelRouteCatalog) ProcessSwaggerJSON() (string, error) {
	mdc.RLock()
	defer mdc.RUnlock()
	allURLs := map[string]map[string]interface{}{}

	for _, mri := range mdc.routes {
		for url, data := range mri.Data {
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
