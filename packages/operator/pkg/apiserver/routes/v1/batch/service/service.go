/*
 *
 *     Copyright 2021 EPAM Systems
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 */

package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/batch"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/httputil"
	logutils "github.com/odahu/odahu-flow/packages/operator/pkg/utils/log"
	"net/http"
	"reflect"
)

const (
	GetURL    = "/batch/service/:id"
	ListURL   = "/batch/service"
	PostURL   = "/batch/service"
	PutURL    = "/batch/service"
	DeleteURL = "/batch/service/:id"
	idParam   = "id"
)

var (
	fieldsCache = map[string]int{}
)

func init() {
	elem := reflect.TypeOf(&batch.InferenceServiceFilter{}).Elem()
	for i := 0; i < elem.NumField(); i++ {
		tagName := elem.Field(i).Tag.Get(batch.TagKey)

		fieldsCache[tagName] = i
	}
}

type Service interface {
	Create(ctx context.Context, bis *batch.InferenceService) (err error)
	Update(ctx context.Context, id string, bis *batch.InferenceService) (err error)
	Delete(ctx context.Context, id string) (err error)
	Get(ctx context.Context, id string) (res batch.InferenceService, err error)
	List(ctx context.Context, options ...filter.ListOption) (res []batch.InferenceService, err error)
}

type controller struct {
	service Service
}

func SetupRoutes(routes gin.IRoutes, service Service) {
	controller := controller{service: service}
	routes.GET(GetURL, controller.Get)
	routes.GET(ListURL, controller.List)
	routes.POST(PostURL, controller.Post)
	routes.PUT(PutURL, controller.Put)
	routes.DELETE(DeleteURL, controller.Delete)
}

// @Summary Get an InferenceService
// @Description Get an InferenceService by id
// @Tags Batch
// @Name id
// @Accept  json
// @Produce  json
// @Param id path string true "InferenceService id"
// @Success 200 {object} batch.InferenceService
// @Failure 404 {object} httputil.HTTPResult
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/batch/service/{id} [get]
func (cr *controller) Get(c *gin.Context) {
	serviceID := c.Param(idParam)

	ctx := c.Request.Context()
	log := logutils.FromContext(ctx)

	service, err := cr.service.Get(ctx, serviceID)
	if err != nil {
		code := errors.CalculateHTTPStatusCode(err)
		if code == http.StatusInternalServerError {
			log.Error(err, fmt.Sprintf("Retrieving %s InferenceService", serviceID))
		}
		c.AbortWithStatusJSON(code, httputil.HTTPResult{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

// @Summary Create an InferenceService
// @Description Create an InferenceService
// @Tags Batch
// @Accept  json
// @Produce  json
// @Param service body batch.InferenceService true "InferenceService". Only `id` and `spec` are taken into account
// @Success 201 {object} batch.InferenceService
// @Failure 404 {object} httputil.HTTPResult
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/batch/service [post]
func (cr *controller) Post(c *gin.Context) {

	var service batch.InferenceService

	ctx := c.Request.Context()
	log := logutils.FromContext(ctx)

	if err := c.ShouldBindJSON(&service); err != nil {
		log.Error(err, "JSON binding of the InferenceService is failed")
		c.AbortWithStatusJSON(http.StatusBadRequest, httputil.HTTPResult{Message: err.Error()})

		return
	}

	err := cr.service.Create(ctx, &service)
	if err != nil {
		code := errors.CalculateHTTPStatusCode(err)
		if code == http.StatusInternalServerError {
			log.Error(err, fmt.Sprintf("Creating %s InferenceService", service.ID))
		}
		c.AbortWithStatusJSON(code, httputil.HTTPResult{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, service)
}

// @Summary Update an InferenceService
// @Description Update an InferenceService
// @Tags Batch
// @Accept  json
// @Produce  json
// @Param service body batch.InferenceService true "InferenceService". Only `id` and `spec` are taken into account
// @Success 200 {object} batch.InferenceService
// @Failure 404 {object} httputil.HTTPResult
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/batch/service [put]
func (cr *controller) Put(c *gin.Context) {

	var service batch.InferenceService

	ctx := c.Request.Context()
	log := logutils.FromContext(ctx)

	if err := c.ShouldBindJSON(&service); err != nil {
		log.Error(err, "JSON binding of the InferenceService is failed")
		c.AbortWithStatusJSON(http.StatusBadRequest, httputil.HTTPResult{Message: err.Error()})

		return
	}

	err := cr.service.Update(ctx, service.ID, &service)
	if err != nil {
		code := errors.CalculateHTTPStatusCode(err)
		if code == http.StatusInternalServerError {
			log.Error(err, fmt.Sprintf("Updating %s InferenceService", service.ID))
		}
		c.AbortWithStatusJSON(code, httputil.HTTPResult{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

// @Summary Delete an InferenceService
// @Description Delete an InferenceService
// @Tags Batch
// @Accept  json
// @Produce  json
// @Param id path string true "InferenceService id"
// @Success 200 {object} httputil.HTTPResult
// @Failure 404 {object} httputil.HTTPResult
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/batch/service/{id} [delete]
func (cr *controller) Delete(c *gin.Context) {

	serviceID := c.Param(idParam)

	ctx := c.Request.Context()
	log := logutils.FromContext(ctx)

	err := cr.service.Delete(ctx, serviceID)
	if err != nil {
		code := errors.CalculateHTTPStatusCode(err)
		if code == http.StatusInternalServerError {
			log.Error(err, fmt.Sprintf("Deleting %s InferenceService", serviceID))
		}
		c.AbortWithStatusJSON(code, httputil.HTTPResult{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, httputil.HTTPResult{Message: fmt.Sprintf("Inference Service %s was deleted", serviceID)})
}

// @Summary List an InferenceService
// @Description List an InferenceService
// @Tags Batch
// @Accept  json
// @Produce  json
// @Param size query int false "Number of entities in a response"
// @Param page query int false "Number of a page"
// @Success 200 {array} batch.InferenceService
// @Failure 404 {object} httputil.HTTPResult
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/batch/service [get]
func (cr *controller) List(c *gin.Context) {

	ctx := c.Request.Context()
	log := logutils.FromContext(ctx)

	f := &batch.InferenceServiceFilter{}
	size, page, err := routes.URLParamsToFilter(c, f, fieldsCache)
	if err != nil {
		log.Error(err, "Malformed url parameters of inference service request")
		c.AbortWithStatusJSON(http.StatusBadRequest, httputil.HTTPResult{Message: err.Error()})

		return
	}

	res, err := cr.service.List(ctx, filter.Size(size), filter.Page(page))
	if err != nil {
		code := errors.CalculateHTTPStatusCode(err)
		if code == http.StatusInternalServerError {
			log.Error(err, "Listing InferenceService")
		}
		c.AbortWithStatusJSON(code, httputil.HTTPResult{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
