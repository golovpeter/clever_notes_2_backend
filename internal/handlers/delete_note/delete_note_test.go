package delete_note

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/golovpeter/clever_notes_2/internal/common/make_error_response"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"testing"
)

type Suite struct {
	suite.Suite

	ctl *gomock.Controller

	responseWriterMock *MockResponseWriter
	databaseMock       *MockDatabase

	handler *deleteNoteHandler
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

	s.handler = NewDeleteNoteHandler(s.databaseMock)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

var testJson = `
	{
		"note_id": 10
	}
`

func (s *Suite) Test_ServerHTTP_Success() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Method: http.MethodPost,
		Body:   body,
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

	s.databaseMock.EXPECT().Exec(gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)

	errorMessage := make_error_response.ErrorMessage{
		ErrorCode:    "0",
		ErrorMessage: "note successful deleted",
	}

	marshaled, err := json.Marshal(errorMessage)

	s.responseWriterMock.EXPECT().Write(marshaled).Times(1).Return(len(marshaled), nil)

	s.handler.ServeHTTP(s.responseWriterMock, &request)

}
