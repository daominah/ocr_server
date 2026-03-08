package httpsvr

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/daominah/ocr_server/pkg/driver/ocr"
)

func Base64Handler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Base64    string `json:"base64"`
		Languages string `json:"languages"`
		Whitelist string `json:"whitelist"`
		// ErodeRadius thins bold characters by removing foreground pixels
		// near edges, helping separate overlapping glyphs in captcha images.
		// 0 means no erosion. Typical values: 1-3.
		ErodeRadius int `json:"erode_radius"`
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
		writeError(w, http.StatusBadRequest, fmt.Errorf("base64 DecodeString: %v", err))
		return
	}
	if isJpg {
		b, err = convertJpgToPng(b)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("convertJpgToPng: %v", err))
			return
		}
	}

	preview := body.Base64
	if len(preview) > 32 {
		preview = preview[:32]
	}
	logPrefix := fmt.Sprintf("base64ed [%.1f kB, %v]", float64(len(b))/1024, preview)

	recognizeAndRespond(w, b, ocr.Params{
		Languages:   body.Languages,
		Whitelist:   body.Whitelist,
		ErodeRadius: body.ErodeRadius,
	}, logPrefix)
}

// recognizeAndRespond runs optional erosion, OCR, and writes the JSON response.
func recognizeAndRespond(w http.ResponseWriter, imageBytes []byte, p ocr.Params, logPrefix string) {
	var err error
	const maxErodeRadius = 10
	if p.ErodeRadius > maxErodeRadius {
		p.ErodeRadius = maxErodeRadius
	}
	if p.ErodeRadius > 0 {
		imageBytes, err = ocr.ErodeImage(imageBytes, p.ErodeRadius)
		if err != nil {
			writeError(w, http.StatusBadRequest, fmt.Errorf("erode image: %v", err))
			return
		}
	}

	var langs []string
	if p.Languages != "" {
		langs = strings.Split(p.Languages, ",")
	}

	var result string
	escapeHTML := true
	switch p.Format {
	case "hocr":
		result, err = ocr.RecognizeHOCR(imageBytes, langs, p.Whitelist)
		escapeHTML = false
	default:
		result, err = ocr.Recognize(imageBytes, langs, p.Whitelist)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("response of %v: %v", logPrefix, result)
	writeJSON(w, http.StatusOK, map[string]interface{}{"result": result}, escapeHTML)
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

func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 256 MiB max in-memory buffering per request; excess spills to disk.
	// Each concurrent upload can use up to this much RAM.
	r.ParseMultipartForm(256 << 20)
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

	erodeRadius, _ := strconv.Atoi(r.FormValue("erode_radius"))

	logPrefix := fmt.Sprintf("file [%.1f kB, %v]", float64(fileHeader.Size)/1024, fileHeader.Filename)

	recognizeAndRespond(w, imageBytes, ocr.Params{
		Languages:   r.FormValue("languages"),
		Whitelist:   r.FormValue("whitelist"),
		Format:      r.FormValue("format"),
		ErodeRadius: erodeRadius,
	}, logPrefix)
}
