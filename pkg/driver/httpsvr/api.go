package httpsvr

import "net/http"

func RegisterAPI(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/status", logMiddleware(Status))
	mux.HandleFunc("POST /api/base64", logMiddleware(Base64))
	mux.HandleFunc("POST /api/file", logMiddleware(FileUpload))
	mux.HandleFunc("/", logMiddleware(Index))
}
