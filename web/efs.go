package web

import "embed"

//go:embed "js"
var JS embed.FS

//go:embed "css"
var CSS embed.FS

//go:embed "fonts"
var Fonts embed.FS
