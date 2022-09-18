package servestatic

import (
	"net/http"
	"os"
	"path/filepath"
)

type serveStaticHandler struct {
	staticPath string
}

const indexFileName = "index.html"

func NewServeStaticHandler(staticPath string) *serveStaticHandler {
	return &serveStaticHandler{staticPath: staticPath}
}

// ServeHTTP возвращает файл из staticPath, если такой файл есть, иначе отдает index.html
func (s *serveStaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(s.staticPath, r.URL.Path)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(s.staticPath, indexFileName))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.FileServer(http.Dir(s.staticPath)).ServeHTTP(w, r)
}
