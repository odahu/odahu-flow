package deployment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	"net/http"
	"strconv"
)

func ValidateAndParseCursor(c *gin.Context, cursor *int) (ok bool) {
	var err error
	cursorParam := c.Query("cursor")
	if cursorParam != "" {
		*cursor, err = strconv.Atoi(cursorParam)
		numErr, isNumErr := err.(*strconv.NumError)
		if err != nil && isNumErr {
			text := "Incorrect \"cursor\" query parameter value: %v. Integer expected. Details: %v"
			c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{
				Message: fmt.Sprintf(text, cursorParam, numErr),
			})
			return false
		} else if err != nil && !isNumErr {
			c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})
			return false
		}
	}
	return true
}