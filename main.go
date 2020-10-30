package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/daominah/ocr_server/controllers"
	"github.com/daominah/ocr_server/filters"
	mylog "github.com/mywrap/log"
	"github.com/otiai10/marmoset"
)

var logger *log.Logger

func main() {

	marmoset.LoadViews("./app/views")

	r := marmoset.NewRouter()
	// API
	r.GET("/status", controllers.Status)
	r.POST("/base64", controllers.Base64)
	r.POST("/file", controllers.FileUpload)
	// Sample Page
	r.GET("/", controllers.Index)
	r.Static("/assets", "./app/assets")

	logger = log.New(os.Stdout, fmt.Sprintf("[%s] ", "ocrserver"), 0)
	r.Apply(&filters.LogFilter{Logger: logger})

	port := os.Getenv("PORT")
	if port == "" {
		mylog.Printf("env `PORT` undefined, listen on default :35735\n")
		port = ":35735"
	}
	if !strings.Contains(port, ":") {
		port = ":" + port
	}
	mylog.Printf("listening on port http://127.0.0.1%s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		mylog.Println(err)
	}
}
