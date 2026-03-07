package httpsvr

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/daominah/ocr_server/pkg/driver/ocr"
	"github.com/otiai10/gosseract/v2"
)

var (
	indexTmpl  *template.Template
	staticServer http.Handler
)

// InitViews parses HTML templates and prepares static file serving from files.
func InitViews(files fs.FS) error {
	staticServer = http.FileServer(http.FS(files))
	var err error
	indexTmpl, err = template.ParseFS(files, "index.html")
	return err
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		staticServer.ServeHTTP(w, r)
		return
	}
	if indexTmpl == nil {
		http.Error(w, "views not initialized", http.StatusInternalServerError)
		return
	}
	indexTmpl.Execute(w, map[string]interface{}{
		"AppName": "Optical Character Recognition server",
	})
}

func Status(w http.ResponseWriter, r *http.Request) {
	langs, err := gosseract.GetAvailableLanguages()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":             "Hello",
		"tesseract_version":   ocr.Version(),
		"tesseract_languages": langs,
	}, true)
}

func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next(w, r)
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}, escapeHTML bool) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(escapeHTML)
	enc.Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()}, true)
}
