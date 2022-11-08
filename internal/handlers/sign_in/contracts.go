package sign_in

import (
	"database/sql"
	"net/http"
)

//go:generate mockgen -source=contracts.go -destination=mocks.go -package=$GOPACKAGE

type Database interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...any) (sql.Result, error)
}

type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}
