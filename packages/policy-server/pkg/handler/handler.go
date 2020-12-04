package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type PolicyStore interface {
	Create(policy OPAPolicy) error
	Delete(ID string) error
	Update(ID string) error
}

type RegoSyntaxError struct {
	error
	Line int
	Column int
}

type RegoSyntaxErrorSet struct {
	error
	Errors []RegoSyntaxError
}

type ErrorWithHTTPCode interface {
	error
	GetHTTPCode() int
}


type OPAPolicy struct {
	ID string `json:"id"`
	RegoPolicy []byte `json:"regoPolicy"`
	Labels map[string]string `json:"labels"`
}

type PolicyHandler struct {
	PolicyStore    PolicyStore
}

func validate(policy OPAPolicy) error {
	return nil
}

func (h *PolicyHandler) PostPolicyHandler(c *gin.Context) {
	var policy OPAPolicy
	if err := c.ShouldBindJSON(&policy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate(policy); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.PolicyStore.Create(policy); err != nil {

		switch err := err.(type) {
		case RegoSyntaxErrorSet:
			syntaxErrors := []struct {
				Line int
			}{{}}
			for _, err := range err.Errors {
				syntaxErrors = append(syntaxErrors, struct{ Line int }{Line: err.Line})
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
				"syntaxErrors": syntaxErrors,
			})
		case ErrorWithHTTPCode:
			c.JSON(err.GetHTTPCode(), gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Something goes wrong... We already work on it"})
		}
		nErr, ok := err.(ErrorWithHTTPCode)
		if ok {
			c.JSON(nErr.GetHTTPCode(), gin.H{"error": nErr.Error()})
			return
		}
	}

}