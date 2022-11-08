package sign_in

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/golovpeter/clever_notes_2/internal/common/hasher"
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

	handler *signInHandler
}

func (s *Suite) SetupTest() {
	s.ctl = gomock.NewController(s.T())

	s.responseWriterMock = NewMockResponseWriter(s.ctl)
	s.databaseMock = NewMockDatabase(s.ctl)

	generateJWT = func(username string) (string, error) {
		return "some jwt", nil
	}

	generateRefreshJWT = func() (string, error) {
		return "some refresh jwt", nil
	}

	s.handler = NewSignInHandler(s.databaseMock)
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
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destBool := dest.(*bool)
			*destBool = true
			return nil
		},
	)

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destUser := dest.(*User)
			*destUser = User{
				User_id:  0,
				Username: "testuser",
				Password: hasher.GeneratePasswordHash("123qweasd"),
			}
			return nil
		})

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)

	accessToken, err := generateJWT("testuser")
	require.NoError(s.T(), err)

	refreshToken, err := generateRefreshJWT()
	require.NoError(s.T(), err)

	s.databaseMock.EXPECT().Exec(gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)

	response := SignInOut{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	marshaled, _ := json.Marshal(response)

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

func (s *Suite) Test_ServeHTTP_NotRegistered() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	errorMessage := make_error_response.ErrorMessage{
		ErrorCode:    "1",
		ErrorMessage: "The user is not registered!",
	}

	marshaled, err := json.Marshal(errorMessage)
	require.NoError(s.T(), err)

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).Return(nil)
	s.responseWriterMock.EXPECT().WriteHeader(http.StatusBadRequest)
	s.responseWriterMock.EXPECT().Write(marshaled).Return(len(marshaled), nil)

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

func (s *Suite) Test_ServeHTTP_WrongGetElementExist() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).Return(errors.New("some error"))
	s.responseWriterMock.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServeHTTP_IncorrectCredentials() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	errorMessage := make_error_response.ErrorMessage{
		ErrorCode:    "1",
		ErrorMessage: "Incorrect password!",
	}

	marshaled, err := json.Marshal(errorMessage)
	require.NoError(s.T(), err)

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destBool := dest.(*bool)
			*destBool = true
			return nil
		},
	)

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).Return(nil)
	s.responseWriterMock.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)
	s.responseWriterMock.EXPECT().Write(marshaled).Times(1).Return(len(marshaled), nil)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServeHTTP_WrongGetTokenExist() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destBool := dest.(*bool)
			*destBool = true
			return nil
		},
	)

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destUser := dest.(*User)
			*destUser = User{
				User_id:  0,
				Username: "testuser",
				Password: hasher.GeneratePasswordHash("123qweasd"),
			}
			return nil
		})

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(errors.New("some error"))
	s.responseWriterMock.EXPECT().WriteHeader(http.StatusInternalServerError)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServeHTTP_WrongTokenExist() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destBool := dest.(*bool)
			*destBool = true
			return nil
		},
	)

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destUser := dest.(*User)
			*destUser = User{
				User_id:  0,
				Username: "testuser",
				Password: hasher.GeneratePasswordHash("123qweasd"),
			}
			return nil
		})

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destBool := dest.(*bool)
			*destBool = true
			return nil
		},
	)

	s.databaseMock.EXPECT().Exec(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("some err_or"))
	s.responseWriterMock.EXPECT().WriteHeader(http.StatusInternalServerError)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServeHTTP_WrongInsertToken() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destBool := dest.(*bool)
			*destBool = true
			return nil
		},
	)

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destUser := dest.(*User)
			*destUser = User{
				User_id:  0,
				Username: "testuser",
				Password: hasher.GeneratePasswordHash("123qweasd"),
			}
			return nil
		})

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)

	_, err := generateJWT("testuser")
	require.NoError(s.T(), err)

	_, err = generateRefreshJWT()
	require.NoError(s.T(), err)

	s.databaseMock.EXPECT().Exec(gomock.Any(), gomock.Any()).Times(1).Return(nil, errors.New("some error"))
	s.responseWriterMock.EXPECT().WriteHeader(http.StatusInternalServerError)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}

func (s *Suite) Test_ServeHTTP_WrongMarshall() {
	body := io.NopCloser(bytes.NewReader([]byte(testJson)))

	request := http.Request{
		Body:   body,
		Method: http.MethodPost,
	}

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destBool := dest.(*bool)
			*destBool = true
			return nil
		},
	)

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), "testuser").Times(1).DoAndReturn(
		func(dest interface{}, query string, args ...interface{}) error {
			destUser := dest.(*User)
			*destUser = User{
				User_id:  0,
				Username: "testuser",
				Password: hasher.GeneratePasswordHash("123qweasd"),
			}
			return nil
		})

	s.databaseMock.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(nil)

	accessToken, err := generateJWT("testuser")
	require.NoError(s.T(), err)

	refreshToken, err := generateRefreshJWT()
	require.NoError(s.T(), err)

	s.databaseMock.EXPECT().Exec(gomock.Any(), gomock.Any()).Times(1).Return(nil, nil)

	wrongResponse := SignInOut{
		AccessToken:  "wrong token",
		RefreshToken: "wrong refresh token",
	}

	response := SignInOut{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	marshaled, _ := json.Marshal(response)
	wrongMarshaled, _ := json.Marshal(wrongResponse)

	s.responseWriterMock.EXPECT().Write(marshaled).Return(len(wrongMarshaled), nil)
	s.responseWriterMock.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)

	s.handler.ServeHTTP(s.responseWriterMock, &request)
}
