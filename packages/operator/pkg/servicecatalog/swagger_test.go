package servicecatalog_test

import (
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/servicecatalog"
	"github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/routes"
	. "github.com/onsi/gomega"
	"github.com/rakyll/statik/fs"
	"github.com/stretchr/testify/suite"
)

const (
	pathPrefix            = "/custom-prefix/example"
	testSwaggerDefinition = "some swagger definition"
)

type SwaggerRouteSuite struct {
	suite.Suite
	g      *GomegaWithT
	server *gin.Engine
}

func (s *SwaggerRouteSuite) SetupSuite() {
	staticFS, err := fs.New()
	if err != nil {
		s.T().Fatal(err)
	}

	s.server = gin.Default()
	rg := s.server.Group(pathPrefix)
	servicecatalog.SetUpSwagger(rg, staticFS, func() (s string, e error) {
		return testSwaggerDefinition, nil
	})
}

func (s *SwaggerRouteSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func TestSwaggerRouteSuite(t *testing.T) {
	suite.Run(t, new(SwaggerRouteSuite))
}

func (s *SwaggerRouteSuite) TestSwaggerRootPage() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, path.Join(pathPrefix, "/swagger/index.html"), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(w.Body.String()).Should(ContainSubstring("SwaggerUIBundle"))
	// Verify that index.html contains custom url.
	// The url is part of our contract.
	// A swagger definition will be available there.
	s.g.Expect(w.Body.String()).Should(ContainSubstring("\"./data.json\""))
}

func (s *SwaggerRouteSuite) verifyMimeType(url, mimeType string) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(w.Header().Get(routes.ContentTypeHeaderKey)).Should(ContainSubstring(mimeType))
}

func (s *SwaggerRouteSuite) TestSwaggerMimeType() {
	s.verifyMimeType(path.Join(pathPrefix, "/swagger/index.html"), "text/html")
	s.verifyMimeType(path.Join(pathPrefix, "/swagger/swagger-ui.css"), "text/css")
	s.verifyMimeType(path.Join(pathPrefix, "/swagger/swagger-ui.js"), "application/javascript")
}

func (s *SwaggerRouteSuite) TestAPISwaggerDefinition() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, path.Join(pathPrefix, "/swagger/data.json"), nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	// Verify random data from API swagger definition
	s.g.Expect(w.Body.String()).Should(Equal(testSwaggerDefinition))
}
