package get_all_notes

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
)

type Suite struct {
	suite.Suite

	ctl *gomock.Controller

	responseWriterMock *MockResponseWriter
	databaseMock       *MockDatabase

	handler *getAllNotesHandel
}

func (s *Suite) SetupTest() {
	s.ctl = gomock.NewController(s.T())

	s.responseWriterMock = NewMockResponseWriter(s.ctl)
	s.databaseMock = NewMockDatabase(s.ctl)

	parseAuthHeader = func(w http.ResponseWriter, r *http.Request) (string, error) {
		return "some token", nil
	}

	validateToken = func(accessToken string) error {
		return nil
	}

	s.handler = NewGetAllNotesHandler(s.databaseMock)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_ServerHTTP_Success() {

	request := http.Request{
		Method: http.MethodPost,
	}

	accessToken, err := parseAuthHeader(s.responseWriterMock, &request)
	require.NoError(s.T(), err)

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), accessToken).Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destBool := dest.(*bool)
			*destBool = true
			return nil
		},
	)

	err = validateToken(accessToken)
	require.NoError(s.T(), err)

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), accessToken).Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destInt := dest.(*int)
			*destInt = 10
			return nil
		},
	)

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destInt := dest.(*int)
			*destInt = 10
			return nil
		},
	)

	s.databaseMock.EXPECT().Select(gomock.Any(), gomock.Any(), 10).Return(nil)

	notes := make([]Note, 0)

	response, err := json.Marshal(GetAllNotesOut{Notes: notes})
	require.NoError(s.T(), err)

	s.responseWriterMock.EXPECT().Write(response).Times(1).Return(len(response), nil)

	s.handler.ServeHTTP(s.responseWriterMock, &request)

}
