package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/otiai10/gosseract/v2"
)

func TestReadCaptcha(t *testing.T) {
	for _, test := range []struct {
		imagePath string
		expected  string
	}{
		{"test/captcha0.png", "JSXJ"},
		//{"test/captcha1.png", "RAF J"},
		{"test/captcha2.png", "CUXJ"},
		{"test/captcha3.png", "OJPJ"},
		{"test/img0.png", "ocrserver"},
		{"test/img1.png", "B-Trees"},
	} {
		client := gosseract.NewClient()
		client.Languages = []string{"eng"}
		err := client.SetConfigFile("tesseract.cfg")
		if err != nil {
			t.Fatal(err)
		}
		defer client.Close()

		image, err := ioutil.ReadFile(test.imagePath)
		if err != nil {
			t.Fatal(err)
		}

		tmpFile, err := ioutil.TempFile("", "tmp_orc_test_")
		if err != nil {
			t.Fatal(err)
		}

		defer func() { tmpFile.Close(); os.Remove(tmpFile.Name()) }()
		tmpFile.Write(image)

		client.SetImage(tmpFile.Name())
		text, err := client.Text()
		if err != nil {
			t.Fatal(err)
		}
		if text != test.expected {
			t.Errorf("error bad read %v, expected: %v on %v",
				text, test.expected, test.imagePath)
		}
	}
}
