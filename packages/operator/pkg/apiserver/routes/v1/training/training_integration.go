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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	httputil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/httputil"
	"net/http"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var logTI = logf.Log.WithName("training-integration-controller")

const (
	GetTrainingIntegrationURL    = "/training-integration/:id"
	GetAllTrainingIntegrationURL = "/training-integration"
	CreateTrainingIntegrationURL = "/training-integration"
	UpdateTrainingIntegrationURL = "/training-integration"
	DeleteTrainingIntegrationURL = "/training-integration/:id"
	IDTiURLParam                 = "id"
)

var (
	emptyCache = map[string]int{}
)

type trainingIntegrationService interface {
	GetTrainingIntegration(name string) (*training.TrainingIntegration, error)
	GetTrainingIntegrationList(options ...filter.ListOption) ([]training.TrainingIntegration, error)
	CreateTrainingIntegration(md *training.TrainingIntegration) error
	UpdateTrainingIntegration(md *training.TrainingIntegration) error
	DeleteTrainingIntegration(name string) error
}

type TrainingIntegrationController struct { //nolint
	service   trainingIntegrationService
	validator *TiValidator
}

// @Summary Get a TrainingIntegration
// @Description Get a TrainingIntegration by id
// @Tags TrainingIntegration
// @Name id
// @Accept  json
// @Produce  json
// @Param id path string true "TrainingIntegration id"
// @Success 200 {object} training.TrainingIntegration
// @Failure 404 {object} httputil.HTTPResult
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/training-integration/{id} [get]
func (tic *TrainingIntegrationController) getTrainingIntegration(c *gin.Context) {
	tiID := c.Param(IDTiURLParam)

	ti, err := tic.service.GetTrainingIntegration(tiID)
	if err != nil {
		logTI.Error(err, fmt.Sprintf("Retrieving %s training integration", tiID))
		c.AbortWithStatusJSON(errors.CalculateHTTPStatusCode(err), httputil.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, ti)
}

// @Summary Get list of TrainingIntegrations
// @Description Get list of TrainingIntegrations
// @Tags TrainingIntegration
// @Accept  json
// @Produce  json
// @Param size path int false "Number of entities in a response"
// @Param page path int false "Number of a page"
// @Success 200 {array} training.TrainingIntegration
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/training-integration [get]
func (tic *TrainingIntegrationController) getAllTrainingIntegrations(c *gin.Context) {
	size, page, err := routes.URLParamsToFilter(c, nil, emptyCache)
	if err != nil {
		logTI.Error(err, "Malformed url parameters of training integration request")
		c.AbortWithStatusJSON(http.StatusBadRequest, httputil.HTTPResult{Message: err.Error()})

		return
	}

	tiList, err := tic.service.GetTrainingIntegrationList(
		filter.Size(size),
		filter.Page(page),
	)
	if err != nil {
		logTI.Error(err, "Retrieving list of training integrations")
		c.AbortWithStatusJSON(errors.CalculateHTTPStatusCode(err), httputil.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, &tiList)
}

// @Summary Create a TrainingIntegration
// @Description Create a TrainingIntegration. Results is created TrainingIntegration.
// @Param ti body training.TrainingIntegration true "Create a TrainingIntegration"
// @Tags TrainingIntegration
// @Accept  json
// @Produce  json
// @Success 201 {object} training.TrainingIntegration
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/training-integration [post]
func (tic *TrainingIntegrationController) createTrainingIntegration(c *gin.Context) {
	var ti training.TrainingIntegration

	if err := c.ShouldBindJSON(&ti); err != nil {
		logTI.Error(err, "JSON binding of training integration is failed")
		c.AbortWithStatusJSON(http.StatusBadRequest, httputil.HTTPResult{Message: err.Error()})

		return
	}

	if err := tic.validator.ValidatesAndSetDefaults(&ti); err != nil {
		logMT.Error(err, fmt.Sprintf("Validation of the training integration is failed: %v", ti))
		c.AbortWithStatusJSON(http.StatusBadRequest, httputil.HTTPResult{Message: err.Error()})

		return
	}

	if err := tic.service.CreateTrainingIntegration(&ti); err != nil {
		logTI.Error(err, fmt.Sprintf("Creation of training integration: %v", ti))
		c.AbortWithStatusJSON(errors.CalculateHTTPStatusCode(err), httputil.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusCreated, ti)
}

// @Summary Update a TrainingIntegration
// @Description Update a TrainingIntegration. Results is updated TrainingIntegration.
// @Param ti body training.TrainingIntegration true "Update a TrainingIntegration"
// @Tags TrainingIntegration
// @Accept  json
// @Produce  json
// @Success 200 {object} training.TrainingIntegration
// @Failure 404 {object} httputil.HTTPResult
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/training-integration [put]
func (tic *TrainingIntegrationController) updateTrainingIntegration(c *gin.Context) {
	var ti training.TrainingIntegration

	if err := c.ShouldBindJSON(&ti); err != nil {
		logTI.Error(err, "JSON binding of training integration is failed")
		c.AbortWithStatusJSON(http.StatusBadRequest, httputil.HTTPResult{Message: err.Error()})

		return
	}

	if err := tic.validator.ValidatesAndSetDefaults(&ti); err != nil {
		logMT.Error(err, fmt.Sprintf("Validation of the training integration is failed: %v", ti))
		c.AbortWithStatusJSON(http.StatusBadRequest, httputil.HTTPResult{Message: err.Error()})

		return
	}

	if err := tic.service.UpdateTrainingIntegration(&ti); err != nil {
		logTI.Error(err, fmt.Sprintf("Update of training integration: %v", ti))
		c.AbortWithStatusJSON(errors.CalculateHTTPStatusCode(err), httputil.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, ti)
}

// @Summary Delete a TrainingIntegration
// @Description Delete a TrainingIntegration by id
// @Tags TrainingIntegration
// @Name id
// @Accept  json
// @Produce  json
// @Param id path string true "TrainingIntegration id"
// @Success 200 {object} httputil.HTTPResult
// @Failure 404 {object} httputil.HTTPResult
// @Failure 400 {object} httputil.HTTPResult
// @Router /api/v1/training-integration/{id} [delete]
func (tic *TrainingIntegrationController) deleteTrainingIntegration(c *gin.Context) {
	tiID := c.Param(IDTiURLParam)

	if err := tic.service.DeleteTrainingIntegration(tiID); err != nil {
		logTI.Error(err, fmt.Sprintf("Deletion of %s training integration is failed", tiID))
		c.AbortWithStatusJSON(errors.CalculateHTTPStatusCode(err), httputil.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, httputil.HTTPResult{Message: fmt.Sprintf("TrainingIntegration %s was deleted", tiID)})
}
