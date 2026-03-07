package web

import "embed"

//go:embed *.html *.css *.js *.ico
var Files embed.FS
