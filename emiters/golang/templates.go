package golang

import _ "embed"

//go:embed templates/model.tmpl
var ModelTemplate string

//go:embed templates/model_field.tmpl
var ModelFieldTemplate string

//go:embed templates/create_constructor.tmpl
var CreateConstructorTemplate string

//go:embed templates/base.tmpl
var BaseTemplate string

//go:embed templates/validator.tmpl
var ValidatorTemplate string
