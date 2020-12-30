package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/mywrap/gofast"
	"github.com/mywrap/log"
	"github.com/otiai10/gosseract/v2"
	"github.com/otiai10/marmoset"
)

// Base64 ...
func Base64(w http.ResponseWriter, r *http.Request) {

	render := marmoset.Render(w, true)

	var body = new(struct {
		Base64    string `json:"base64"`
		Trim      string `json:"trim"`
		Languages string `json:"languages"`
		Whitelist string `json:"whitelist"`
	})

	err := json.NewDecoder(r.Body).Decode(body)
	if err != nil {
		render.JSON(http.StatusBadRequest, err)
		return
	}

	tempfile, err := ioutil.TempFile("", "ocrserver"+"-")
	if err != nil {
		render.JSON(http.StatusInternalServerError, err)
		return
	}
	defer func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}()

	if len(body.Base64) == 0 {
		render.JSON(http.StatusBadRequest, fmt.Errorf("base64 string required"))
		return
	}
	isJpg := false
	if strings.HasPrefix(body.Base64, `data:image/jpeg;base64,`) {
		isJpg = true
	}
	body.Base64 = strings.TrimPrefix(body.Base64, `data:image/png;base64,`)
	body.Base64 = strings.TrimPrefix(body.Base64, `data:image/jpeg;base64,`)
	b, err := base64.StdEncoding.DecodeString(body.Base64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("error base64 DecodeString: %v", err)))
		return
	}
	if isJpg {
		b, err = ConvertJpgToPng(b)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("error ConvertJpgToPng: %v", err)))
			return
		}
	}
	tempfile.Write(b)

	client := gosseract.NewClient()
	defer client.Close()
	err = client.SetConfigFile("./tesseract.cfg")
	if err != nil {
		log.Printf("error tesseract SetConfigFile: %v", err)
	}
	client.Languages = []string{"eng"}
	if body.Languages != "" {
		client.Languages = strings.Split(body.Languages, ",")
	}
	client.SetImage(tempfile.Name())
	if body.Whitelist != "" {
		client.SetWhitelist(body.Whitelist)
	}

	text, err := client.Text()
	if err != nil {
		render.JSON(http.StatusInternalServerError, err)
		return
	}

	ret := strings.Trim(text, body.Trim)
	log.Debugf("response of base64ed [%.1f kB, %v]: %v", float64(len(b))/1024,
		body.Base64[:gofast.MinInts(32, len(body.Base64))], ret)
	render.JSON(http.StatusOK, map[string]interface{}{"result": ret})
}

func ConvertJpgToPng(imageBytes []byte) ([]byte, error) {
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
