package weeny

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type URLParseSuite struct {
	suite.Suite
}

func TestURLParse(t *testing.T) {
	suite.Run(t, new(URLParseSuite))
}

func (s *URLParseSuite) SetupSuite() {
}

func (s *URLParseSuite) TearDownSuite() {
}
