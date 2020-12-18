package servicecatalogroutes

import (
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils/swagger"
	"github.com/swaggo/swag"
	"net/http"
)

func SetUpSwagger(rg *gin.RouterGroup, apiStaticFS http.FileSystem) {
	rg.GET("/swagger/*any", swagger.Handler(apiStaticFS, swag.ReadDoc))
}