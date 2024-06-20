package weeny

import (
	"testing"

	"github.com/k0kubun/pp/v3"
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

func (s *URLParseSuite) TestGetDomain() {
	tests := []struct {
		uri  string
		want string
	}{
		{
			uri:  `https://crc.wintalent.cn/wt/CRC/web/index/CompCRCPagerecruit_Social`,
			want: "crc.wintalent.cn",
		},
	}
	for _, tt := range tests {
		got := DomainFromURL(tt.uri)

		s.Equal(tt.want, got)
	}
}

func (s *URLParseSuite) TestURL() {
	tests := []struct {
		uri  string
		want string
	}{
		{
			uri:  `https://crc.wintalent.cn/wt/CRC/web/templet1000/index/corpwebPosition1000CRC!getOnePosition?postIdEnc=bbc6f69131d2c881fd4a2dffde1d1a86&brandCode=1&recruitType=2&lanType=1&showComp=true`,
			want: "crc.wintalent.cn",
		},
	}
	for _, tt := range tests {
		uu, err := ParseURL(tt.uri)
		s.Nil(err)

		pp.Println(uu)
	}
}
