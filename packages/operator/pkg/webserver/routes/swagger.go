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

package routes

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/odahu/odahu-flow/packages/operator/docs" //nolint
	_ "github.com/odahu/odahu-flow/packages/operator/pkg/static"
	"github.com/swaggo/swag"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var logS = logf.Log.WithName("swagger-controller")

const (
	// This path is hardcoded in the packages/operator/static/swagger/index.html file
	pathToSwaggerDefinition = "data.json"
	// It is a directory in the virtual filesystem.
	// Dir location is there packages/operator/static/swagger
	swaggerDir = "/swagger"
	// octet-stream is a binary format
	defaultMimeType      = "application/octet-stream"
	ContentTypeHeaderKey = "Content-Type"
)

// Error messages
const (
	WritingDataErrorMessage              = "Writing swagger data is failed"
	ReadingFromVFSErrorMessage           = "Error while reading swagger files from virtual file system"
	ReadingSwaggerDefinitionErrorMessage = "Reading of a swagger definition is failed"
)

func SetUpSwagger(rg *gin.RouterGroup, apiStaticFS http.FileSystem) {
	rg.GET("/swagger/*any", SwaggerHandler(apiStaticFS, swag.ReadDoc))
}

// It returns a v2 or v3 swagger definition as a string.
// Example of the definition is https://petstore.swagger.io/v2/swagger.json
type SwaggerDefinitionReader func() (string, error)

func SwaggerHandler(apiStaticFS http.FileSystem, swaggerDefReader SwaggerDefinitionReader) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileName := path.Base(c.Request.URL.Path)

		switch fileName {
		case pathToSwaggerDefinition:
			processSwaggerDefinition(swaggerDefReader, c)
		default:
			processSwaggerFiles(apiStaticFS, fileName, c)
		}
	}
}

func processSwaggerFiles(staticFS http.FileSystem, fileName string, c *gin.Context) {
	file, err := staticFS.Open(path.Join(swaggerDir, fileName))
	if os.ErrNotExist == err {
		c.AbortWithStatusJSON(http.StatusNotFound, HTTPResult{
			Message: fmt.Sprintf("not found: %s", c.Request.URL),
		})

		return
	}

	if err != nil {
		logS.Error(err, ReadingFromVFSErrorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, HTTPResult{Message: ReadingFromVFSErrorMessage})

		return
	}

	mimeType := mime.TypeByExtension(filepath.Ext(fileName))
	if mimeType == "" {
		mimeType = defaultMimeType
	}

	c.Header(ContentTypeHeaderKey, mimeType)
	c.Status(http.StatusOK)

	_, err = io.Copy(c.Writer, file)
	if err != nil {
		logS.Error(err, WritingDataErrorMessage)

		c.AbortWithStatusJSON(http.StatusInternalServerError, HTTPResult{Message: WritingDataErrorMessage})
	}
}

func processSwaggerDefinition(swaggerDefReader SwaggerDefinitionReader, c *gin.Context) {
	swaggerDef, err := swaggerDefReader()
	if err != nil {
		logS.Error(err, ReadingSwaggerDefinitionErrorMessage)

		c.AbortWithStatusJSON(http.StatusInternalServerError, HTTPResult{Message: ReadingSwaggerDefinitionErrorMessage})

		return
	}

	_, err = c.Writer.Write([]byte(swaggerDef))
	if err != nil {
		logS.Error(err, WritingDataErrorMessage)

		c.AbortWithStatusJSON(http.StatusInternalServerError, HTTPResult{Message: WritingDataErrorMessage})
	} else {
		c.Status(http.StatusOK)
	}
}
