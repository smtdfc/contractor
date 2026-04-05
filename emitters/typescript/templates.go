package typescript

import _ "embed"

//go:embed templates/runtime.tmpl
var RuntimeTemplate string

//go:embed templates/model.tmpl
var ModelTemplate string

//go:embed templates/error.tmpl
var ErrorTemplate string

//go:embed templates/constructor.tmpl
var ConstructorTemplate string

//go:embed templates/validator_method.tmpl
var ValidatorMethodTemplate string

//go:embed templates/rest.tmpl
var RestTemplate string

//go:embed templates/base.tmpl
var BaseTemplate string
