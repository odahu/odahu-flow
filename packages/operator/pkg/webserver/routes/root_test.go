package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/odahu/odahu-flow/packages/operator/pkg/webserver/routes"
	. "github.com/onsi/gomega"
	"github.com/rakyll/statik/fs"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

const (
	// Title from packages/operator/pkg/static/index.html file
	rootPageTitle    = "Odahu-flow API Gateway web console"
	partOfCLICommand = "odahuflowctl login --url "
	testLink1Name    = "test-link-1"
	testLink1URL     = "http://test.org/test1"
	testLink2Name    = "test-link-2"
	testLink2URL     = "http://test.org/test2"
)

type RootRouteSuite struct {
	suite.Suite
	g      *GomegaWithT
	server *gin.Engine
}

func (s *RootRouteSuite) SetupSuite() {
	var links = []interface{}{
		map[interface{}]interface{}{
			"name": testLink1Name,
			"url":  testLink1URL,
		},
		map[interface{}]interface{}{
			"name": testLink2Name,
			"url":  testLink2URL,
		},
	}
	viper.Set("common.external_urls", links)

	s.server = gin.Default()
	routerGroup := s.server.Group("")

	staticFS, err := fs.New()
	if err != nil {
		s.T().Fatal(err)
	}

	err = routes.SetUpIndexPage(routerGroup, staticFS)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *RootRouteSuite) SetupTest() {
	s.g = NewGomegaWithT(s.T())
}

func TestRootRouteSuite(t *testing.T) {
	suite.Run(t, new(RootRouteSuite))
}

func (s *RootRouteSuite) TestRootPage() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(w.Body.String()).Should(ContainSubstring(rootPageTitle))
	s.g.Expect(w.Body.String()).Should(ContainSubstring(partOfCLICommand))
}

func (s *RootRouteSuite) TestRootPageToken() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	s.g.Expect(err).Should(BeNil())

	testToken := "test-token"
	req.Header.Add(routes.AuthHeaderName, testToken)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(w.Body.String()).Should(ContainSubstring(testToken))
}

func (s *RootRouteSuite) TestRootPageLinks() {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	s.g.Expect(err).NotTo(HaveOccurred())
	s.server.ServeHTTP(w, req)

	s.g.Expect(w.Code).Should(Equal(http.StatusOK))
	s.g.Expect(w.Body.String()).Should(ContainSubstring(testLink1Name))
	s.g.Expect(w.Body.String()).Should(ContainSubstring(testLink1URL))
	s.g.Expect(w.Body.String()).Should(ContainSubstring(testLink2Name))
	s.g.Expect(w.Body.String()).Should(ContainSubstring(testLink2URL))
}
