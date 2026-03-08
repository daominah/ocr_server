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
// This file path is absolute path, work in the Docker container.
const cfgFile = "/tesseract.cfg"

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

// Params holds OCR options shared across handlers.
type Params struct {
	// Languages is a comma-separated list of Tesseract language codes
	// (e.g. "eng", "vie", "chi_sim"). When empty, defaults to "eng"
	// with captcha-optimized config (no dictionary penalties).
	Languages string
	// Whitelist limits recognized characters to this set
	// (e.g. "abcdefghijklmnopqrstuvwxyz0123456789").
	Whitelist string
	// ErodeRadius thins bold characters by removing foreground pixels
	// near edges, helping separate overlapping glyphs in captcha images.
	// 0 means no erosion. Typical values: 1-3.
	ErodeRadius int
	// Format selects the output format: "" for plain text, "hocr" for
	// hOCR (HTML with bounding boxes and confidence scores).
	Format string
}

// Version returns the installed Tesseract version string.
func Version() string {
	client := gosseract.NewClient()
	defer client.Close()
	return client.Version()
}
