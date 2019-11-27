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
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apis/configuration"

	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var logR = logf.Log.WithName("root-controller")

const (
	// Path of the file is packages/operator/static/index.html file
	pathToRootFile       = "/index.html"
	AuthHeaderName       = "X-Jwt"
	RootFileTemplateName = "root_index_template"
)

// Error messages
const (
	TemplateParsingErrorMessage = "Failed template parsing of root page"
)

type RootFileDataTemplate struct {
	Token string
	Links []configuration.ExternalUrl
}

func SetUpIndexPage(routerGroup *gin.RouterGroup, apiStaticFS http.FileSystem) (err error) {
	rootFile, err := apiStaticFS.Open(pathToRootFile)
	if err != nil {
		return err
	}
	defer func() {
		err = rootFile.Close()
	}()

	rootFileData, err := ioutil.ReadAll(rootFile)
	if err != nil {
		return err
	}

	rootFileTemplate, err := template.New(RootFileTemplateName).Parse(string(rootFileData))
	if err != nil {
		return err
	}

	externalUrls := configuration.ExportExternalUrlsFromConfig()

	routerGroup.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
		err := rootFileTemplate.Execute(c.Writer, RootFileDataTemplate{
			Token: c.GetHeader(AuthHeaderName),
			Links: externalUrls,
		})

		if err != nil {
			logR.Error(err, TemplateParsingErrorMessage)
			c.AbortWithStatusJSON(http.StatusInternalServerError, TemplateParsingErrorMessage)
			return
		}
	})

	return nil
}
