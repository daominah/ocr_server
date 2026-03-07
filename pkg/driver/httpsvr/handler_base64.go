package httpsvr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	"log"

	"github.com/daominah/ocr_server/pkg/driver/ocr"
)

func Base64(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Base64    string `json:"base64"`
		Trim      string `json:"trim"`
		Languages string `json:"languages"`
		Whitelist string `json:"whitelist"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if len(body.Base64) == 0 {
		writeError(w, http.StatusBadRequest, fmt.Errorf("base64 string required"))
		return
	}

	isJpg := strings.HasPrefix(body.Base64, `data:image/jpeg;base64,`)
	body.Base64 = strings.TrimPrefix(body.Base64, `data:image/png;base64,`)
	body.Base64 = strings.TrimPrefix(body.Base64, `data:image/jpeg;base64,`)

	b, err := base64.StdEncoding.DecodeString(body.Base64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error base64 DecodeString: %v", err)))
		return
	}
	if isJpg {
		b, err = convertJpgToPng(b)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("error convertJpgToPng: %v", err)))
			return
		}
	}

	var langs []string
	if body.Languages != "" {
		langs = strings.Split(body.Languages, ",")
	}
	text, err := ocr.Recognize(b, langs, body.Whitelist)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	ret := strings.Trim(text, body.Trim)
	preview := body.Base64
	if len(preview) > 32 {
		preview = preview[:32]
	}
	log.Printf("response of base64ed [%.1f kB, %v]: %v",
		float64(len(b))/1024, preview, ret)
	writeJSON(w, http.StatusOK, map[string]interface{}{"result": ret}, true)
}

func convertJpgToPng(imageBytes []byte) ([]byte, error) {
	imgType := http.DetectContentType(imageBytes)
	switch imgType {
	case "image/png":
		return imageBytes, nil
	case "image/jpeg":
		img, err := jpeg.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, fmt.Errorf("decode jpeg: %v", err)
		}
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			return nil, fmt.Errorf("encode png: %v", err)
		}
		return buf.Bytes(), nil
	default:
		return nil, fmt.Errorf("unsupported image type %v", imgType)
	}
}
