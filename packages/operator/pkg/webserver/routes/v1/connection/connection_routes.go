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

package connection

import (
	"fmt"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/connection"
	conn_repository "github.com/odahu/odahu-flow/packages/operator/pkg/repository/connection"
	conn_service "github.com/odahu/odahu-flow/packages/operator/pkg/service/connection"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var logC = logf.Log.WithName("connection-controller")

const (
	GetConnectionURL           = "/connection/:id"
	GetDecryptedConnectionURL  = "/connection/:id/decrypted"
	GetAllConnectionURL        = "/connection"
	CreateConnectionURL        = "/connection"
	UpdateConnectionURL        = "/connection"
	DeleteConnectionURL        = "/connection/:id"
	IDConnURLParam             = "id"
	ConnDecryptTokenQueryParam = "token"
)

var (
	fieldsCache = map[string]int{}
)

func init() {
	elem := reflect.TypeOf(&conn_repository.Filter{}).Elem()
	for i := 0; i < elem.NumField(); i++ {
		tagName := elem.Field(i).Tag.Get(conn_repository.TagKey)

		fieldsCache[tagName] = i
	}
}

type controller struct {
	connService conn_service.Service
	validator   *ConnValidator
}

func ConfigureRoutes(
	routeGroup *gin.RouterGroup,
	connService conn_service.Service,
	keyEvaluator PublicKeyEvaluator,
	connectionConfig config.ConnectionConfig,
) {
	controller := &controller{
		connService: connService,
		validator:   NewConnValidator(keyEvaluator),
	}
	routeGroup = routeGroup.Group("", routes.DisableAPIMiddleware(connectionConfig.Enabled))

	routeGroup.GET(GetConnectionURL, controller.getConnection)
	routeGroup.GET(GetDecryptedConnectionURL, controller.getDecryptedConnection)
	routeGroup.GET(GetAllConnectionURL, controller.getAllConnections)
	routeGroup.POST(CreateConnectionURL, controller.createConnection)
	routeGroup.PUT(UpdateConnectionURL, controller.updateConnection)
	routeGroup.DELETE(DeleteConnectionURL, controller.deleteConnection)
}

// @Summary Get a Connection
// @Description Get a Connection by id
// @Tags Connection
// @Name id
// @Accept  json
// @Produce  json
// @Param id path string true "Connection id"
// @Success 200 {object} connection.Connection
// @Failure 404 {object} routes.HTTPResult
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/connection/{id} [get]
func (cc *controller) getConnection(c *gin.Context) {
	connID := c.Param(IDConnURLParam)

	conn, err := cc.connService.GetConnection(connID, true)
	if err != nil {
		logC.Error(err, fmt.Sprintf("Retrieving %s connection", connID))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, conn)
}

// @Summary Get a decrypted Connection
// @Description Get a decrypted Connection by id
// @Tags Connection
// @Name id
// @Accept  json
// @Produce  json
// @Param id path string true "Connection id"
// @Success 200 {object} connection.Connection
// @Failure 404 {object} routes.HTTPResult
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/connection/{id}/decrypted [get]
func (cc *controller) getDecryptedConnection(c *gin.Context) {
	connID := c.Param(IDConnURLParam)

	conn, err := cc.connService.GetConnection(connID, false)
	if err != nil {
		logC.Error(err, fmt.Sprintf("Retrieving %s connection", connID))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, conn)
}

// @Summary Get list of Connections
// @Description Get list of Connections
// @Tags Connection
// @Accept  json
// @Produce  json
// @Param type path string false "Toolchain"
// @Param size path int false "Number of entities in a response"
// @Param page path int false "Number of a page"
// @Success 200 {array} connection.Connection
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/connection [get]
func (cc *controller) getAllConnections(c *gin.Context) {
	filter := &conn_repository.Filter{}
	size, page, err := routes.URLParamsToFilter(c, filter, fieldsCache)
	if err != nil {
		logC.Error(err, "Malformed url parameters of connection request")
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	connList, err := cc.connService.GetConnectionList(
		conn_repository.ListFilter(filter),
		conn_repository.Size(size),
		conn_repository.Page(page),
	)
	if err != nil {
		logC.Error(err, "Retrieving list of connections")
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}
	c.JSON(http.StatusOK, connList)
}

// @Summary Create a Connection
// @Description Create a Connection. Results is created Connection.
// @Param connection body connection.Connection true "Create a Connection"
// @Tags Connection
// @Accept  json
// @Produce  json
// @Success 201 {object} connection.Connection
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/connection [post]
func (cc *controller) createConnection(c *gin.Context) {
	var conn connection.Connection

	if err := c.ShouldBindJSON(&conn); err != nil {
		logC.Error(err, "JSON binding of the connection is failed")
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	if err := cc.validator.ValidatesAndSetDefaults(&conn); err != nil {
		logC.Error(err, fmt.Sprintf("Validation of the connection is failed: %v", conn))
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	var createdConnection *connection.Connection
	var err error
	if createdConnection, err = cc.connService.CreateConnection(conn); err != nil {
		logC.Error(err, fmt.Sprintf("Creation of the connection: %+v", conn))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusCreated, createdConnection)
}

// @Summary Update a Connection
// @Description Update a Connection. Results is updated Connection.
// @Param connection body connection.Connection true "Update a Connection"
// @Tags Connection
// @Accept  json
// @Produce  json
// @Success 200 {object} routes.HTTPResult
// @Failure 404 {object} routes.HTTPResult
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/connection [put]
func (cc *controller) updateConnection(c *gin.Context) {
	var conn connection.Connection

	if err := c.ShouldBindJSON(&conn); err != nil {
		logC.Error(err, "JSON binding of the connection is failed")
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	if err := cc.validator.ValidatesAndSetDefaults(&conn); err != nil {
		logC.Error(err, fmt.Sprintf("Validation of the connection is failed: %v", conn))
		c.AbortWithStatusJSON(http.StatusBadRequest, routes.HTTPResult{Message: err.Error()})

		return
	}

	updatedConnection, err := cc.connService.UpdateConnection(conn)
	if err != nil {
		logC.Error(err, fmt.Sprintf("Update of the connection: %+v", conn))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, updatedConnection)
}

// @Summary Delete a Connection
// @Description Delete a Connection by id
// @Tags Connection
// @Name id
// @Accept  json
// @Produce  json
// @Param id path string true "Connection id"
// @Success 200 {object} routes.HTTPResult
// @Failure 404 {object} routes.HTTPResult
// @Failure 400 {object} routes.HTTPResult
// @Router /api/v1/connection/{id} [delete]
func (cc *controller) deleteConnection(c *gin.Context) {
	connID := c.Param(IDConnURLParam)

	if err := cc.connService.DeleteConnection(connID); err != nil {
		logC.Error(err, fmt.Sprintf("Deletion of %s connection is failed", connID))
		c.AbortWithStatusJSON(routes.CalculateHTTPStatusCode(err), routes.HTTPResult{Message: err.Error()})

		return
	}

	c.JSON(http.StatusOK, routes.HTTPResult{Message: fmt.Sprintf("Connection %s was deleted", connID)})
}
