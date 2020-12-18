package servicecatalogroutes_test

import (
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/model"
	"github.com/odahu/odahu-flow/packages/operator/pkg/servicecatalog/servicecatalogroutes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/odahu/odahu-flow/packages/operator/pkg/servicecatalog"
)

func TestDeployedModelHandler (t *testing.T) {
	log, err := zap.NewDevelopment()
	assert.NoError(t, err)


	catalog := servicecatalog.NewModelRouteCatalog(log.Sugar())

	raw, err := ioutil.ReadFile("testdata/swagger.json")
	assert.NoError(t, err)
	assert.NoError(t, catalog.CreateOrUpdate(servicecatalog.Route{
		ID:     "simple",
		Prefix: "/model/simple",
		IsDefault: true,
		Model: model.DeployedModel{
			DeploymentID: "simple-model",
			ServedModel: model.ServedModel{
				Metadata: model.Metadata{},
				MLServer: "Triton",
				Swagger:  model.Swagger2{Raw: raw},
			},
		},
	}))

	engine := gin.Default()
	servicecatalogroutes.SetupDeployedModelRoute(engine.Group(""), catalog.GetDeployedModel)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/model-info/simple-model", nil)

	assert.NoError(t, err)

	engine.ServeHTTP(w, req)

	assert.Equal(t, w.Code, http.StatusOK)

}