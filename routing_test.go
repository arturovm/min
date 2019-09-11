package min_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/arturovm/min"
	"github.com/arturovm/min/mocks"
)

type RoutingTestSuite struct {
	suite.Suite
	g              *min.Group
	h              *mocks.Handler
	registeredPath string
}

func (s *RoutingTestSuite) SetupTest() {
	s.h = new(mocks.Handler)
	s.g = min.New(s.h).NewGroup("/sub")
	s.registeredPath = "/sub/path"
}

func (s *RoutingTestSuite) on(method string) {
	s.h.On("Handle", method, s.registeredPath, nil)
}

func (s *RoutingTestSuite) TestGet() {
	s.on(http.MethodGet)
	s.g.Get("/path", nil)
	s.h.AssertExpectations(s.T())
}

func (s *RoutingTestSuite) TestPost() {
	s.on(http.MethodPost)
	s.g.Post("/path", nil)
	s.h.AssertExpectations(s.T())
}

func (s *RoutingTestSuite) TestPut() {
	s.on(http.MethodPut)
	s.g.Put("/path", nil)
	s.h.AssertExpectations(s.T())
}

func (s *RoutingTestSuite) TestPatch() {
	s.on(http.MethodPatch)
	s.g.Patch("/path", nil)
	s.h.AssertExpectations(s.T())
}

func (s *RoutingTestSuite) TestDelete() {
	s.on(http.MethodDelete)
	s.g.Delete("/path", nil)
	s.h.AssertExpectations(s.T())
}

func TestRouting(t *testing.T) {
	suite.Run(t, new(RoutingTestSuite))
}
