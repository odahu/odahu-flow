/*
 *
 *     Copyright 2021 EPAM Systems
 *
 *     Licensed under the Apache License, Version 2.0 (the "License");
 *     you may not use this file except in compliance with the License.
 *     You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *     Unless required by applicable law or agreed to in writing, software
 *     distributed under the License is distributed on an "AS IS" BASIS,
 *     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *     See the License for the specific language governing permissions and
 *     limitations under the License.
 */

package log

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// The key type is unexported to prevent collisions with context keys defined in
// other package
type key int

const (
	loggerKey key = 0
	loggerName string = "api-request-handler"
)

// If x-request-id header is found then bind it to logger as request-id
// Otherwise generate unique
func getRequestID(c *gin.Context) string {
	if id := c.Request.Header.Get("x-request-id"); len(id) > 0 {
		return id
	}
	return uuid.New().String()
}

func SetupLogBindingMiddleware(r gin.IRoutes) {
	r.Use(func(c *gin.Context) {
		ctx := c.Request.Context()
		var log = logf.Log.WithName(loggerName)
		log = log.WithValues("request-id", getRequestID(c),
			"path", c.Request.URL.Path, "method", c.Request.Method)
		ctx = context.WithValue(ctx, loggerKey, log)
		c.Request = c.Request.WithContext(ctx)
		log.Info("Start handling request")
		c.Next()
		log.Info("Response", "Code", c.Writer.Status())
	})
}

// If logr.Logger stored at loggerKey in context.Context then
// returns logger with appropriate bound `values`.
// Otherwise returns logr.Logger without Request specific bound values
func FromContext(ctx context.Context) logr.Logger {
	v := ctx.Value(loggerKey)
	log, ok := v.(logr.Logger)
	if ok {
		return log
	}
	return logf.Log.WithName(loggerName)
}
