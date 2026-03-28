package golang

import _ "embed"

//go:embed templates/model.tmpl
var ModelTemplate string

//go:embed templates/model_field.tmpl
var ModelFieldTemplate string
