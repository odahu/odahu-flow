package routes

import (
	"github.com/gin-gonic/gin"
	model_types "github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"
	odahu_errors "github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	"net/http"
)

const GetDeployedModelInfoURL = "/model-info/:id"

type GetDeployedModelFunc func(deploymentID string) (model_types.DeployedModel, error)

type DeployedModelHandler struct {
	GetDeployedModel GetDeployedModelFunc
}

type DeployedModelHandlerParams struct {
	ID string `uri:"id" binding:"required"`
}
func (t *DeployedModelHandler) Handle(c *gin.Context) {

	var params DeployedModelHandlerParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(400, gin.H{"msg": "not correct path"})
	}

	model, err := t.GetDeployedModel(params.ID)
	if err != nil {
		c.AbortWithStatusJSON(odahu_errors.CalculateHTTPStatusCode(err), gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusOK, model)
}

func SetupDeployedModelRoute(rg *gin.RouterGroup, getter GetDeployedModelFunc) {
	handler := DeployedModelHandler{GetDeployedModel: getter}
	rg.GET(GetDeployedModelInfoURL, handler.Handle)
}