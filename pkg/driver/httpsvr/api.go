package httpsvr

import "net/http"

func RegisterAPI(mux *http.ServeMux) {
	mux.HandleFunc("/", logMiddleware(Index))
	mux.HandleFunc("GET /api/status", logMiddleware(StatusHandler))

	// OCR endpoints:

	mux.HandleFunc("POST /api/base64", logMiddleware(Base64Handler))
	mux.HandleFunc("POST /api/file", logMiddleware(FileUploadHandler))
}
