package deployment

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/errors"
	httputil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/httputil"
	"net/http"
	"strconv"
)

func ValidateAndParseCursor(c *gin.Context, cursor *int) (err error) {
	cursorParam := c.Query("cursor")
	if cursorParam != "" {
		*cursor, err = strconv.Atoi(cursorParam)
		numErr, isNumErr := err.(*strconv.NumError)
		if err != nil && isNumErr {
			text := "Incorrect \"cursor\" query parameter value: %v. Integer expected. Details: %v"
			c.AbortWithStatusJSON(http.StatusBadRequest, httputil.HTTPResult{
				Message: fmt.Sprintf(text, cursorParam, numErr),
			})
		} else if err != nil && !isNumErr {
			c.AbortWithStatusJSON(errors.CalculateHTTPStatusCode(err), httputil.HTTPResult{Message: err.Error()})
		}
	}
	return
}