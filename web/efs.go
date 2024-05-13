package web

import "embed"

//go:embed "js"
var Files embed.FS

//go:embed "css"
var CSS embed.FS
