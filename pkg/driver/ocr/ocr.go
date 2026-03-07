// Package ocr wraps the gosseract CGo bindings for Tesseract OCR.
// It requires Tesseract headers and libraries at build time; go vet/build
// will fail on systems without them (e.g. Windows without Tesseract installed).
package ocr

import (
	"log"

	"github.com/otiai10/gosseract/v2"
)

// cfgFile disables Tesseract's language model penalties, suitable for
// captcha-like images where dictionary constraints hurt accuracy.
const cfgFile = "tesseract.cfg"

func newClient(imageBytes []byte, languages []string, whitelist string) (*gosseract.Client, error) {
	client := gosseract.NewClient()
	if len(languages) == 0 {
		// No language specified: use 0-penalty config for captcha-style OCR.
		if err := client.SetConfigFile(cfgFile); err != nil {
			log.Printf("tesseract SetConfigFile: %v", err)
		}
		languages = []string{"eng"}
	}
	// Languages specified: rely on Tesseract's default penalties for better
	// natural-language accuracy.
	client.Languages = languages
	if whitelist != "" {
		client.SetWhitelist(whitelist)
	}
	if err := client.SetImageFromBytes(imageBytes); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}

// Recognize runs OCR on imageBytes and returns plain text.
func Recognize(imageBytes []byte, languages []string, whitelist string) (string, error) {
	client, err := newClient(imageBytes, languages, whitelist)
	if err != nil {
		return "", err
	}
	defer client.Close()
	return client.Text()
}

// RecognizeHOCR runs OCR on imageBytes and returns hOCR formatted output.
func RecognizeHOCR(imageBytes []byte, languages []string, whitelist string) (string, error) {
	client, err := newClient(imageBytes, languages, whitelist)
	if err != nil {
		return "", err
	}
	defer client.Close()
	return client.HOCRText()
}

// Version returns the installed Tesseract version string.
func Version() string {
	client := gosseract.NewClient()
	defer client.Close()
	return client.Version()
}
