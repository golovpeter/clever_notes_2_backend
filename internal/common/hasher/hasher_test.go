//package hasher
//
//import (
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/suite"
//	"testing"
//)

//func TestGeneratePasswordHash(t *testing.T) {
//	cases := []struct {
//		name     string
//		password string
//		expected string
//	}{
//		{
//			"ok", "123qweasd", "57ba172a6be125cca2f449826f9980ca",
//		},
//	}
//
//	for _, d := range cases {
//		t.Run(d.name, func(t *testing.T) {
//			result := GeneratePasswordHash(d.password)
//
//			assert.Equal(t, d.expected, result, fmt.Sprintf("Expected %s, got %s", d.expected, result))
//		})
//	}
//}

package hasher

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	ctl *gomock.Controller
}

func (s *Suite) SetupTest() {
	s.ctl = gomock.NewController(s.T())
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_GeneratePasswordHash_Success() {
	result := GeneratePasswordHash("123qweasd")
	expected := "57ba172a6be125cca2f449826f9980ca"

	assert.Equal(
		s.T(),
		expected,
		result,
		fmt.Sprintf("Test fail. Expected %s, got %s", expected, result),
	)
}
