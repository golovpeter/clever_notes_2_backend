package make_error_response

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	ctl *gomock.Controller

	responseWriterMock *MockResponseWriter
}

func (s *Suite) SetupTest() {
	s.ctl = gomock.NewController(s.T())

	s.responseWriterMock = NewMockResponseWriter(s.ctl)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_MakeErrorResponse_Success() {
	errMessage := ErrorMessage{
		ErrorCode:    "1",
		ErrorMessage: "some message",
	}

	marshaled, err := json.Marshal(errMessage)

	require.NoError(s.T(), err)

	s.responseWriterMock.EXPECT().Write(marshaled).Times(1).Return(len(marshaled), nil)

	MakeErrorResponse(s.responseWriterMock, errMessage)
}

func (s *Suite) Test_MakeErrorResponse_WriteError() {
	errMessage := ErrorMessage{
		ErrorCode:    "1",
		ErrorMessage: "some message",
	}

	marshaled, err := json.Marshal(errMessage)

	require.NoError(s.T(), err)

	s.responseWriterMock.EXPECT().Write(marshaled).Times(1).Return(0, errors.New("some error"))
	s.responseWriterMock.EXPECT().WriteHeader(http.StatusInternalServerError).Times(1)

	MakeErrorResponse(s.responseWriterMock, errMessage)
}
