package helpers

type FieldValidateRule struct {
	RuleName string
	Value    string
	Message  string
}

type FieldValidationMetadata struct {
	Name  string
	Rules []FieldValidateRule
}
