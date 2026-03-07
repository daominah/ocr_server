package httpsvr

import (
	"io"
	"net/http"
	"strings"

	"log"

	"github.com/daominah/ocr_server/pkg/driver/ocr"
)

func FileUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	upload, fileHeader, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	defer upload.Close()

	imageBytes, err := io.ReadAll(upload)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	var langs []string
	if l := r.FormValue("languages"); l != "" {
		langs = strings.Split(l, ",")
	}
	whitelist := r.FormValue("whitelist")

	var out string
	escapeHTML := true
	switch r.FormValue("format") {
	case "hocr":
		out, err = ocr.RecognizeHOCR(imageBytes, langs, whitelist)
		escapeHTML = false
	default:
		out, err = ocr.Recognize(imageBytes, langs, whitelist)
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ret := strings.Trim(out, r.FormValue("trim"))
	log.Printf("response of file [%.1f kB, %v]: %v",
		float64(fileHeader.Size)/1024, fileHeader.Filename, ret)
	writeJSON(w, http.StatusOK, map[string]interface{}{"result": ret}, escapeHTML)
}
