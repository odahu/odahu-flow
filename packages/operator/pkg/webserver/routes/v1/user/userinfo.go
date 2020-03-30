/*
 * Copyright 2020 EPAM Systems
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package userinfo

import (
	"fmt"
	request_jwt "github.com/dgrijalva/jwt-go/request"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/user"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/odahu/odahu-flow/packages/operator/pkg/utils"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes"
	"net/http"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"

	"github.com/gin-gonic/gin"
)

// URLs
const (
	GetUserInfoURL = "/user/info"
)

const (
	controllerName = "user_controller"
)

var log = logf.Log.WithName(controllerName)

type controller struct {
	config config.Claims
}

func ConfigureRoutes(routeGroup *gin.RouterGroup, config config.Claims) {
	userController := controller{config: config}

	routeGroup.GET(GetUserInfoURL, userController.getUserInfo)
}

// @Summary Get the user information
// @Description Get the user information(email, name and so on)
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200 {object} user.UserInfo
// @Router /api/v1/user/info [get]
func (cc *controller) getUserInfo(c *gin.Context) {
	token, err := request_jwt.AuthorizationHeaderExtractor.ExtractToken(c.Request)
	if err == request_jwt.ErrNoTokenInRequest {
		c.JSON(http.StatusOK, &user.AnonymousUser)
		return
	} else if err != nil {
		const errorMessage = "Unexpected error during extraction a token from headers"
		log.Error(err, errorMessage)

		c.JSON(http.StatusBadRequest, routes.HTTPResult{
			Message: errorMessage,
		})
		return
	}

	userInfo, err := utils.ExtractUserInfoFromToken(token, cc.config)
	if err != nil {
		c.JSON(http.StatusBadRequest, routes.HTTPResult{
			Message: fmt.Sprintf("Malformed JWT: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}
