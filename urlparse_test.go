package weeny

import (
	"net/url"
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

func (s *URLParseSuite) TestGetHost() {
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
		got := HostFromURL(tt.uri)

		s.Equal(tt.want, got)
	}
}

func (s *URLParseSuite) TestURL() {
	tests := []struct {
		uri  string
		want *url.URL
	}{
		{
			uri: `https://crc.wintalent.cn/wt/CRC/web/templet1000/index/corpwebPosition1000CRC!getOnePosition?postIdEnc=bbc6f69131d2c881fd4a2dffde1d1a86&brandCode=1&recruitType=2&lanType=1&showComp=true`,
			want: &url.URL{
				Scheme:      "https",
				Opaque:      "",
				User:        (*url.Userinfo)(nil),
				Host:        "crc.wintalent.cn",
				Path:        "/wt/CRC/web/templet1000/index/corpwebPosition1000CRC!getOnePosition",
				RawPath:     "/wt/CRC/web/templet1000/index/corpwebPosition1000CRC!getOnePosition",
				OmitHost:    false,
				ForceQuery:  false,
				RawQuery:    "postIdEnc=bbc6f69131d2c881fd4a2dffde1d1a86&brandCode=1&recruitType=2&lanType=1&showComp=true",
				Fragment:    "",
				RawFragment: "",
			},
		},
	}
	for _, tt := range tests {
		uu, err := ParseURL(tt.uri)
		s.Nil(err)

		s.Equal(tt.want, uu)
	}
}
