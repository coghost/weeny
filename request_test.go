package weeny

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RequestSuite struct {
	suite.Suite
}

func TestRequest(t *testing.T) {
	suite.Run(t, new(RequestSuite))
}

func (s *RequestSuite) SetupSuite() {
}

func (s *RequestSuite) TearDownSuite() {
}
