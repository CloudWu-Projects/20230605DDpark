package www

import "embed"

//go:embed *.html
var HtmlFS embed.FS
