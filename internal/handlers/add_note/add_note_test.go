package add_note

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
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

	handler *addNoteHandler
}

func (s *Suite) SetupTest() {
	s.ctl = gomock.NewController(s.T())

	s.responseWriterMock = NewMockResponseWriter(s.ctl)
	s.databaseMock = NewMockDatabase(s.ctl)

	validateToken = func(accessToken string) error {
		return nil
	}

	parseAuthHeader = func(w http.ResponseWriter, r *http.Request) (string, error) {
		return "some token", nil
	}

	s.handler = NewAddNoteHandler(s.databaseMock)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

const testJson = `
	{
		"note_caption": "test_note_caption",
		"note": "test_note"
	}
`

func (s *Suite) Test_ServeHTTP_Success() {
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

	s.databaseMock.EXPECT().Exec(gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)
	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	response := AddNoteOut{
		NoteId: 0,
	}

	marshaled, err := json.Marshal(response)
	require.NoError(s.T(), err)

	s.responseWriterMock.EXPECT().Write(marshaled).Times(1).Return(len(marshaled), nil)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServeHTTP_WrongMethod() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodGet,
	}

	errorMessage := make_error_response.ErrorMessage{
		ErrorCode:    "1",
		ErrorMessage: "Unsupported method",
	}

	marshaled, err := json.Marshal(errorMessage)
	require.NoError(s.T(), err)

	s.responseWriterMock.EXPECT().WriteHeader(http.StatusMethodNotAllowed).Times(1)
	s.responseWriterMock.EXPECT().Write(marshaled).Times(1).Return(len(marshaled), nil)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServeHTTP_WrongParseAuthHeader() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Method: http.MethodPost,
		Body:   body,
	}

	parseAuthHeader = func(w http.ResponseWriter, r *http.Request) (string, error) {
		return "", errors.New("some error")
	}

	_, _ = parseAuthHeader(s.responseWriterMock, &request)

	s.handler.ServeHTTP(s.responseWriterMock, &request)

}

func (s *Suite) Test_ServerHTTP_WrongTokenExistGet() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Method: http.MethodPost,
		Body:   body,
	}

	accessToken, err := parseAuthHeader(s.responseWriterMock, &request)
	require.NoError(s.T(), err)

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), accessToken).Times(1).Return(errors.New("some error"))
	s.responseWriterMock.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServerHTTP_TokenExpired() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Method: http.MethodPost,
		Body:   body,
	}

	validateToken = func(accessToken string) error {
		return jwt.ErrTokenExpired
	}

	errorMessage := make_error_response.ErrorMessage{
		ErrorCode:    "1",
		ErrorMessage: "access token expired",
	}

	marshaled, err := json.Marshal(errorMessage)
	require.NoError(s.T(), err)

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

	s.responseWriterMock.EXPECT().WriteHeader(http.StatusUnauthorized).Times(1)
	s.responseWriterMock.EXPECT().Write(marshaled).Times(1).Return(len(marshaled), nil)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServerHTTP_WrongGetUserID() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Method: http.MethodPost,
		Body:   body,
	}

	errorMessage := make_error_response.ErrorMessage{
		ErrorCode:    "1",
		ErrorMessage: "there is no such user",
	}

	marshaled, err := json.Marshal(errorMessage)
	require.NoError(s.T(), err)

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
			*destInt = 0
			return nil
		})
	s.responseWriterMock.EXPECT().WriteHeader(http.StatusBadRequest).Times(1)
	s.responseWriterMock.EXPECT().Write(marshaled).Times(1).Return(len(marshaled), nil)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}
