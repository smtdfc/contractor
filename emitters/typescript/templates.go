package typescript

import (
	"embed"
	_ "embed"
)

//go:embed templates/*.tmpl
var templateFiles embed.FS
