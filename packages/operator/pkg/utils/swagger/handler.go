package swagger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	httputil "github.com/odahu/odahu-flow/packages/operator/pkg/utils/httputil"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var logS = log.Log.WithName("swagger-controller")

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

// It returns a v2 or v3 swagger definition as a string.
// Example of the definition is https://petstore.swagger.io/v2/swagger.json
type DefinitionReader func() (string, error)

func Handler(apiStaticFS http.FileSystem, swaggerDefReader DefinitionReader) gin.HandlerFunc {
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
		c.AbortWithStatusJSON(http.StatusNotFound, httputil.HTTPResult{
			Message: fmt.Sprintf("not found: %s", c.Request.URL),
		})

		return
	}

	if err != nil {
		logS.Error(err, ReadingFromVFSErrorMessage)
		c.AbortWithStatusJSON(http.StatusInternalServerError, httputil.HTTPResult{Message: ReadingFromVFSErrorMessage})

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

		c.AbortWithStatusJSON(http.StatusInternalServerError, httputil.HTTPResult{Message: WritingDataErrorMessage})
	}
}

func processSwaggerDefinition(swaggerDefReader DefinitionReader, c *gin.Context) {
	swaggerDef, err := swaggerDefReader()
	if err != nil {
		logS.Error(err, ReadingSwaggerDefinitionErrorMessage)

		c.AbortWithStatusJSON(http.StatusInternalServerError,
			httputil.HTTPResult{Message: ReadingSwaggerDefinitionErrorMessage})

		return
	}

	_, err = c.Writer.Write([]byte(swaggerDef))
	if err != nil {
		logS.Error(err, WritingDataErrorMessage)

		c.AbortWithStatusJSON(http.StatusInternalServerError, httputil.HTTPResult{Message: WritingDataErrorMessage})
	} else {
		c.Status(http.StatusOK)
	}
}
