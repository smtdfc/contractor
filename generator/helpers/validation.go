package helpers

var SupportedValidators = map[string]bool{
	"IsEmail":       true,
	"Max":           true,
	"Min":           true,
	"IsInt":         true,
	"IsFloat":       true,
	"IsBoolean":     true,
	"IsString":      true,
	"IsDateString":  true,
	"IsUUID":        true,
	"IsUrl":         true,
	"IsArray":       true,
	"MinLength":     true,
	"MaxLength":     true,
	"Length":        true,
	"IsNotEmpty":    true,
	"ArrayMinSize":  true,
	"ArrayMaxSize":  true,
	"ArrayLength":   true,
	"IsPhoneNumber": true,
}

type FieldValidateRule struct {
	RuleName string
	Value    string
	Message  string
}

type FieldValidationMetadata struct {
	Name  string
	Rules []FieldValidateRule
}
