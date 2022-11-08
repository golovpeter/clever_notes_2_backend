package make_error_response

import "net/http"

//go:generate mockgen -source=contracts.go -destination=mocks.go -package=$GOPACKAGE

type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}
