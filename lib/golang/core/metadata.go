package core

type ModelReflectionMetadata[T any] struct {
	Name      string
	Validator IValidator[T]
}
