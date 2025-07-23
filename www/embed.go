package www

import "embed"

//go:embed www/*
var HtmlFS embed.FS
