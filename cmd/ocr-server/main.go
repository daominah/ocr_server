package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/daominah/ocr_server/pkg/driver/httpsvr"
	webfiles "github.com/daominah/ocr_server/web"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	if err := httpsvr.InitViews(webfiles.Files); err != nil {
		log.Fatalf("init views: %v", err)
	}

	mux := http.NewServeMux()
	httpsvr.RegisterAPI(mux)

	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("env `PORT` undefined, listen on default :35735\n")
		port = ":35735"
	}
	if !strings.Contains(port, ":") {
		port = ":" + port
	}
	log.Printf("listening on port http://127.0.0.1%s", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Println(err)
	}
}
