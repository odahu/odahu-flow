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
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/filter"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/api/v1alpha1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/training"
	mt_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/training"
	mt_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/training"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	kube_client "github.com/odahu/odahu-flow/packages/operator/pkg/kubeclient/trainingclient"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var logMT = logf.Log.WithName("training-controller")

const (
	GetModelTrainingURL        = "/model/training/:id"
	GetAllModelTrainingURL     = "/model/training"
	GetModelTrainingLogsURL    = "/model/training/:id/log"
	CreateModelTrainingURL     = "/model/training"
	UpdateModelTrainingURL     = "/model/training"
	SaveModelTrainingResultURL = "/model/training/:id/result"
	DeleteModelTrainingURL     = "/model/training/:id"
	IDMtURLParam               = "id"
	FollowURLParam             = "follow"
)

var (
	fieldsCache = map[string]int{}
)

func init() {
	elem := reflect.TypeOf(&mt_repository.MTFilter{}).Elem()
	for i := 0; i < elem.NumField(); i++ {
		tagName := elem.Field(i).Tag.Get(mt_repository.TagKey)

		fieldsCache[tagName] = i
	}
}

type ModelTrainingController struct {
	kubeClient   kube_client.Client
	trainService mt_service.Service
	validator    *MtValidator
}

// @Summary Get a Model Training
// @Description Get a Model Training by id
// @Tags Training
// @Name id
// @Accept  json
// @Produce  json
// @Param id path string true "Model Training id"
// @Success 200 {object} training.ModelTraining
// @Failure 404 {object} routes.HTTPResult
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/model/training/{id} [get]
func (mtc *ModelTrainingController) getMT(c *gin.Context) {
	mtID := c.Param(IDMtURLParam)

	mt, err := mtc.trainService.GetModelTraining(c.Request.Context(), mtID)
	if err != nil {
		logMT.Error(err, fmt.Sprintf("Retrieving of %s model training", mtID))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, mt)
}

// @Summary Get list of Model Trainings
// @Description Get list of Model Trainings
// @Tags Training
// @Accept  json
// @Produce  json
// @Param size path int false "Number of entities in a response"
// @Param page path int false "Number of a page"
// @Param model_name path int false "Model name"
// @Param model_version path int false "Model version"
// @Param toolchain path int false "Toolchain name"
// @Success 200 {array} training.ModelTraining
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/model/training [get]
func (mtc *ModelTrainingController) getAllMTs(c *gin.Context) {
	f := &mt_repository.MTFilter{}
	size, page, err := routes.URLParamsToFilter(c, f, fieldsCache)
	if err != nil {
		logMT.Error(err, "Malformed url parameters of model training request")
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	mtList, err := mtc.trainService.GetModelTrainingList(
		c.Request.Context(), filter.ListFilter(f), filter.Size(size), filter.Page(page),
	)
	if err != nil {
		logMT.Error(err, "Retrieving list of model trainings")
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, &mtList)
}

// @Summary Create a Model Training
// @Description Create a Model Training. Results is created Model Training.
// @Param mt body training.ModelTraining true "Create a Model Training"
// @Tags Training
// @Accept  json
// @Produce  json
// @Success 201 {object} training.ModelTraining
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/model/training [post]
func (mtc *ModelTrainingController) createMT(c *gin.Context) {
	var mt training.ModelTraining

	if err := c.ShouldBindJSON(&mt); err != nil {
		logMT.Error(err, "JSON binding of the model training is failed")
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	if err := mtc.validator.ValidatesAndSetDefaults(&mt); err != nil {
		logMT.Error(err, fmt.Sprintf("Validation of the model training is failed: %v", mt))
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	if err := mtc.trainService.CreateModelTraining(c.Request.Context(), &mt); err != nil {
		logMT.Error(err, fmt.Sprintf("Creation of the model training: %v", mt))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusCreated, mt)
}

// @Summary Update a Model Training
// @Description Update a Model Training. Results is updated Model Training.
// @Param mt body training.ModelTraining true "Update a Model Training"
// @Tags Training
// @Accept  json
// @Produce  json
// @Success 200 {object} training.ModelTraining
// @Failure 404 {object} routes.HTTPResult
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/model/training [put]
func (mtc *ModelTrainingController) updateMT(c *gin.Context) {
	var mt training.ModelTraining

	if err := c.ShouldBindJSON(&mt); err != nil {
		logMT.Error(err, fmt.Sprintf("JSON binding of the model training is failed: %v", mt))
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	if err := mtc.validator.ValidatesAndSetDefaults(&mt); err != nil {
		logMT.Error(err, fmt.Sprintf("Creation of the model training: %v", mt))
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	if err := mtc.trainService.UpdateModelTraining(c.Request.Context(), &mt); err != nil {
		logMT.Error(err, fmt.Sprintf("Creation of the model training: %v", mt))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, mt)
}

// @Summary Save a Model Training result
// @Description Save a Model Training by id
// @Tags Training
// @Name id
// @Accept  json
// @Produce  json
// @Param MP body v1alpha1.TrainingResult true "Model Training result"
// @Param id path string true "Model Training id"
// @Success 200 {array} v1alpha1.TrainingResult
// @Failure 404 {object} routes.HTTPResult
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/model/training/{id}/result [put]
func (mtc *ModelTrainingController) saveMTResult(c *gin.Context) {
	mtID := c.Param(IDMtURLParam)
	mtResult := &v1alpha1.TrainingResult{}

	if err := c.ShouldBindJSON(mtResult); err != nil {
		logMT.Error(err, "JSON binding of the model training result is failed")
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	if err := mtc.kubeClient.SaveModelTrainingResult(mtID, mtResult); err != nil {
		logMT.Error(err, fmt.Sprintf("Save the result of the model training: %+v", mtResult))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, mtResult)
}

// @Summary Get a Model Training
// @Description Get a Model Training by id
// @Tags Training
// @Name id
// @Accept  json
// @Produce  json
// @Param id path string true "Model Training id"
// @Success 200 {object} routes.HTTPResult
// @Failure 404 {object} routes.HTTPResult
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/model/training/{id} [delete]
func (mtc *ModelTrainingController) deleteMT(c *gin.Context) {
	mtID := c.Param(IDMtURLParam)

	if err := mtc.trainService.DeleteModelTraining(c.Request.Context(), mtID); err != nil {
		logMT.Error(err, fmt.Sprintf("Deletion of %s model training is failed", mtID))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, routes.HTTPResult{Message: fmt.Sprintf("Model training %s was deleted", mtID)})
}

// @Summary Stream logs from model training pod
// @Description Stream logs from model training pod
// @Tags Training
// @Name id
// @Produce  plain
// @Accept  plain
// @Param follow query bool false "follow logs"
// @Param id path string true "Model Training id"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Router /api/v1/model/training/{id}/log [get]
func (mtc *ModelTrainingController) getModelTrainingLog(c *gin.Context) {
	mtID := c.Param(IDMtURLParam)
	follow := false
	var err error

	urlParameters := c.Request.URL.Query()
	followParam := urlParameters.Get(FollowURLParam)

	if len(followParam) != 0 {
		follow, err = strconv.ParseBool(followParam)
		if err != nil {
			errMessage := fmt.Sprintf("Convert %s to bool", followParam)
			logMT.Error(err, errMessage)
			c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: errMessage})

			return
		}
	}

	if err := mtc.kubeClient.GetModelTrainingLogs(mtID, c.Writer, follow); err != nil {
		logMT.Error(err, fmt.Sprintf("Getting %s model training logs is failed", mtID))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}
}
