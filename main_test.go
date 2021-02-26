package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/otiai10/gosseract/v2"
)

func TestReadCaptcha(t *testing.T) {
	for _, test := range []struct {
		imagePath  string
		limitChars string
		language   string // default eng
		expected   string
	}{
		{imagePath: "test/captcha0.png", expected: "JSXJ",
			limitChars: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{imagePath: "test/captcha1.png", expected: "RAFJ",
			limitChars: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{imagePath: "test/captcha2.png", expected: "CUXJ",
			limitChars: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{imagePath: "test/captcha3.png", expected: "OJPJ",
			limitChars: "ABCDEFGHIJKLMNOPQRSTUVWXYZ"},
		{imagePath: "test/img0.png", expected: "ocrserver"},
		{imagePath: "test/img1.png", expected: "B-Trees"},
		{imagePath: "test/vie0.png", expected: "Đào Thanh Tùng", language: "vie"},
		{imagePath: "test/vie1.png", expected: "Đào Thị Lán", language: "vie"},
		{imagePath: "test/chi0.png", expected: "纤 扬", language: "chi_sim"},
		{imagePath: "test/chi1.png", expected: "鱼", language: "chi_sim"},
		{imagePath: "test/chi2.png", expected: "松", language: "chi_sim"},
		{imagePath: "test/dict0.jpg", expected: "Sài Sơn, Quốc Oai, Hà Nội", language: "vie"},
		{imagePath: "test/dict1.jpg", expected: "ĐÀO THANH TÙNG", language: "vie",
			limitChars: " AÀÁÃẠẢĂẮẰẲẴẶÂẤẦẨẪẬBCDĐEÈÉẸẺẼÊẾỀỂỄỆFGHIÌÍĨỈỊJKLMNOÒÓÕỌỎÔỐỒỔỖỘƠỚỜỞỠỢPQRSTUÙÚŨỤỦƯỨỪỬỮỰVWXYÝỲỴỶỸZ"},
	} {
		client := gosseract.NewClient()
		client.Languages = []string{"eng"}
		err := client.SetConfigFile("tesseract.cfg")
		if err != nil {
			t.Fatal(err)
		}
		if test.language != "" {
			client.Languages = []string{test.language}
		}
		if test.limitChars != "" {
			client.SetWhitelist(test.limitChars)
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
