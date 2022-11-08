package sign_up

import (
	"bytes"
	"encoding/json"
	"errors"
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

	handler *signUpHandler
}

func (s *Suite) SetupTest() {
	s.ctl = gomock.NewController(s.T())

	s.responseWriterMock = NewMockResponseWriter(s.ctl)
	s.databaseMock = NewMockDatabase(s.ctl)

	s.handler = NewSignUpHandler(s.databaseMock)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

const testJson = `
	{
		"username": "testuser",
		"password": "123qweasd"
	}
`

func (s *Suite) Test_ServeHTTP_Success() {
	errorMessage := make_error_response.ErrorMessage{
		ErrorCode:    "0",
		ErrorMessage: "Registration was successful!",
	}

	marshaled, err := json.Marshal(errorMessage)
	require.NoError(s.T(), err)

	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).Return(nil)
	s.databaseMock.EXPECT().Exec(gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)
	s.responseWriterMock.EXPECT().Write(marshaled).Times(1).Return(len(marshaled), nil)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServeHTTP_WrongMethod() {
	errorMessage := make_error_response.ErrorMessage{
		ErrorCode:    "1",
		ErrorMessage: "Unsupported method",
	}

	marshaled, err := json.Marshal(errorMessage)
	require.NoError(s.T(), err)

	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodGet,
	}

	s.responseWriterMock.EXPECT().WriteHeader(http.StatusMethodNotAllowed).Times(1)
	s.responseWriterMock.EXPECT().Write(marshaled).Times(1).Return(len(marshaled), nil)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServeHTTP_WrongInput() {
	const testJson = `
	{
		"username": 12345,
		"password": 12345
	}
`
	errorMessage := make_error_response.ErrorMessage{
		ErrorCode:    "1",
		ErrorMessage: "Incorrect data input",
	}

	marshaled, err := json.Marshal(errorMessage)
	require.NoError(s.T(), err)

	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	s.responseWriterMock.EXPECT().WriteHeader(http.StatusBadRequest).Times(1)
	s.responseWriterMock.EXPECT().Write(marshaled).Times(1).Return(len(marshaled), nil)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServeHTTP_WrongValidate() {
	const testJson = `
	{
		"username": "",
		"password": ""
	}
`

	errorMessage := make_error_response.ErrorMessage{
		ErrorCode:    "1",
		ErrorMessage: "Incorrect data input",
	}

	marshaled, err := json.Marshal(errorMessage)
	require.NoError(s.T(), err)

	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	s.responseWriterMock.EXPECT().WriteHeader(http.StatusBadRequest).Times(1)
	s.responseWriterMock.EXPECT().Write(marshaled).Times(1).Return(len(marshaled), nil)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServeHTTP_WrongGet() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).Return(errors.New("some error"))
	s.responseWriterMock.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServeHTTP_WrongExec() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).Return(nil)
	s.databaseMock.EXPECT().Exec(gomock.Any(), "testuser", gomock.Any()).Times(1).Return(nil, errors.New("some error"))
	s.responseWriterMock.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServeHTTP_ElementExist() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	errorMessage := make_error_response.ErrorMessage{
		ErrorCode:    "1",
		ErrorMessage: "User already registered!",
	}

	marshaled, err := json.Marshal(errorMessage)
	require.NoError(s.T(), err)

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destBool := dest.(*bool)
			*destBool = true
			return nil
		})

	s.responseWriterMock.EXPECT().WriteHeader(http.StatusBadRequest).Times(1)
	s.responseWriterMock.EXPECT().Write(marshaled).Times(1).Return(len(marshaled), nil)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}
