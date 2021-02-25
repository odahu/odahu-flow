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

package job

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
	GetURL    = "/batch/job/:id"
	ListURL   = "/batch/job"
	PostURL   = "/batch/job"
	DeleteURL = "/batch/job/:id"
	idParam   = "id"
)

var (
	fieldsCache = map[string]int{}
)

func init() {
	elem := reflect.TypeOf(&batch.InferenceJobFilter{}).Elem()
	for i := 0; i < elem.NumField(); i++ {
		tagName := elem.Field(i).Tag.Get(batch.JobTagKey)

		fieldsCache[tagName] = i
	}
}


type Service interface {
	Create(ctx context.Context, bij *batch.InferenceJob) (err error)
	SetDeletionMark(ctx context.Context, id string) error
	List(ctx context.Context, options ...filter.ListOption) ([]batch.InferenceJob, error)
	Get(ctx context.Context, id string) (batch.InferenceJob, error)
}

type controller struct {
	service Service
}

func SetupRoutes(routes gin.IRoutes, service Service) {
	c := controller{service: service}
	routes.GET(GetURL, c.Get)
	routes.GET(ListURL, c.List)
	routes.POST(PostURL, c.Post)
	routes.DELETE(DeleteURL, c.Delete)
}

// @Summary Get an InferenceJob
// @Description Get an InferenceJob by id
// @Tags Batch
// @Name id
// @Accept  json
// @Produce  json
// @Param id path string true "InferenceJob id"
// @Success 200 {object} batch.InferenceJob
// @Failure 404 {object} httputil.HTTPResult
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/batch/job/{id} [get]
func (cr *controller) Get(c *gin.Context) {
	id := c.Param(idParam)

	ctx := c.Request.Context()
	log := logutils.FromContext(ctx)

	service, err := cr.service.Get(ctx, id)
	if err != nil {
		code := errors.CalculateHTTPStatusCode(err)
		if code == http.StatusInternalServerError {
			log.Error(err, fmt.Sprintf("Retrieving %s InferenceJob", id))
		}
		c.AbortWithStatusJSON(code, httputil.HTTPResult{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, service)
}

// @Summary Create an InferenceJob
// @Description Create an InferenceJob
// @Tags Batch
// @Accept  json
// @Produce  json
// @Param service body batch.InferenceJob true "InferenceJob". Only `id` and `spec` are taken into account
// @Success 201 {object} batch.InferenceJob
// @Failure 404 {object} httputil.HTTPResult
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/batch/job [post]
func (cr *controller) Post(c *gin.Context) {

	var job batch.InferenceJob

	ctx := c.Request.Context()
	log := logutils.FromContext(ctx)

	if err := c.ShouldBindJSON(&job); err != nil {
		log.Error(err, "JSON binding of the InferenceJob is failed")
		c.AbortWithStatusJSON(http.StatusBadRequest, httputil.HTTPResult{Message: err.Error()})

		return
	}

	err := cr.service.Create(ctx, &job)
	if err != nil {
		code := errors.CalculateHTTPStatusCode(err)
		if code == http.StatusInternalServerError {
			log.Error(err, fmt.Sprintf("Creating %s InferenceJob", job.ID))
		}
		c.AbortWithStatusJSON(code, httputil.HTTPResult{Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, job)
}

// @Summary Delete an InferenceJob
// @Description Delete an InferenceJob
// @Tags Batch
// @Accept  json
// @Produce  json
// @Param id path string true "Model Training id"
// @Success 200 {object} httputil.HTTPResult
// @Failure 404 {object} httputil.HTTPResult
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/batch/job/{id} [delete]
func (cr *controller) Delete(c *gin.Context) {

	id := c.Param(idParam)

	ctx := c.Request.Context()
	log := logutils.FromContext(ctx)

	err := cr.service.SetDeletionMark(ctx, id)
	if err != nil {
		code := errors.CalculateHTTPStatusCode(err)
		if code == http.StatusInternalServerError {
			log.Error(err, fmt.Sprintf("Deleting %s InferenceJob", id))
		}
		c.AbortWithStatusJSON(code, httputil.HTTPResult{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, httputil.HTTPResult{Message: fmt.Sprintf("Inference Service %s was deleted", id)})
}


// @Summary List an InferenceJob
// @Description List an InferenceJob
// @Tags Batch
// @Accept  json
// @Produce  json
// @Param size query int false "Number of entities in a response"
// @Param page query int false "Number of a page"
// @Success 200 {array} batch.InferenceJob
// @Failure 404 {object} httputil.HTTPResult
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/batch/job [get]
func (cr *controller) List(c *gin.Context) {

	ctx := c.Request.Context()
	log := logutils.FromContext(ctx)

	f := &batch.InferenceJobFilter{}
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
			log.Error(err, "Listing InferenceJob")
		}
		c.AbortWithStatusJSON(code, httputil.HTTPResult{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

