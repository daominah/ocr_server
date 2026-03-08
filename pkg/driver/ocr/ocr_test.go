package ocr

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestRecognize(t *testing.T) {
	for _, test := range []struct {
		imagePath string
		whitelist string
		language  string // default eng
		expected  string
	}{
		{imagePath: "testdata/captcha0.png", expected: "JSXJ",
			whitelist: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{imagePath: "testdata/captcha1.png", expected: "RAFJ",
			whitelist: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{imagePath: "testdata/captcha2.png", expected: "CUXJ",
			whitelist: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{imagePath: "testdata/captcha3.png", expected: "OJPJ",
			whitelist: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{imagePath: "testdata/img0.png", expected: "ocrserver"},
		{imagePath: "testdata/img1.png", expected: "B-Trees"},
		{imagePath: "testdata/vie0.png", expected: "Đào Thanh Tùng", language: "vie"},
		{imagePath: "testdata/vie1.png", expected: "Đào Thị Lán", language: "vie"},
		{imagePath: "testdata/chi0.png", expected: "纤 扬", language: "chi_sim"},
		{imagePath: "testdata/chi1.png", expected: "鱼", language: "chi_sim"},
		{imagePath: "testdata/chi2.png", expected: "松", language: "chi_sim"},
		{imagePath: "testdata/dict0.jpg", expected: "Sài Sơn, Quốc Oai, Hà Nội", language: "vie"},
		{imagePath: "testdata/dict1.jpg", expected: "ĐÀO THANH TÙNG", language: "vie",
			whitelist: " AÀÁÃẠẢĂẮẰẲẴẶÂẤẦẨẪẬBCDĐEÈÉẸẺẼÊẾỀỂỄỆFGHIÌÍĨỈỊJKLMNOÒÓÕỌỎÔỐỒỔỖỘƠỚỜỞỠỢPQRSTUÙÚŨỤỦƯỨỪỬỮỰVWXYÝỲỴỶỸZ"},
	} {
		imageBytes, err := os.ReadFile(test.imagePath)
		if err != nil {
			t.Fatal(err)
		}

		var langs []string
		if test.language != "" {
			langs = []string{test.language}
		}

		text, err := Recognize(imageBytes, langs, test.whitelist)
		if err != nil {
			t.Fatal(err)
		}
		if text != test.expected {
			t.Errorf("got %q, want %q for %v", text, test.expected, test.imagePath)
		}
	}
}

// TestRecognizeOverlapChars do OCR for captcha images with overlapping characters.
// Set ERODE_RADIUS env var to test with erosion, typically 1, 2 or 3.
//
//	ERODE_RADIUS=1 go test -v ./pkg/driver/ocr/ -run=TestRecognizeOverlapChars
func TestRecognizeOverlapChars(t *testing.T) {
	erodeRadius, _ := strconv.Atoi(os.Getenv("ERODE_RADIUS"))
	t.Logf("ERODE_RADIUS=%d", erodeRadius)

	paths, err := filepath.Glob("testdata/overlap/*.png")
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range paths {
		imageBytes, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if erodeRadius > 0 {
			imageBytes, err = ErodeImage(imageBytes, erodeRadius)
			if err != nil {
				t.Errorf("%v erode: %v", p, err)
				continue
			}
		}
		text, err := Recognize(imageBytes, nil, "")
		if err != nil {
			t.Errorf("%v: %v", p, err)
			continue
		}
		t.Logf("%v result: %q", filepath.Base(p), text)
	}
}
