package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
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
	body.Base64 = regexp.MustCompile("data:image\\/png;base64,").ReplaceAllString(body.Base64, "")
	b, err := base64.StdEncoding.DecodeString(body.Base64)
	if err != nil {
		render.JSON(http.StatusBadRequest, err)
		return
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
