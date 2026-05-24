# Golang Code Generation

This guide explains exactly what Contractor currently generates for Golang, how to use the generated types, and what to watch out for.

## Scope

Current Golang output focuses on:

- struct types with JSON tags
- struct types with discriminator `__k` embedded
- custom error types
- enum string constants
- generic types
- REST metadata

## Generation Pipeline

1. Parse `.contract` files into AST.
2. Type-check models, references, generics, and annotations.
3. Convert AST to IR.
4. Emit Golang source from IR.

## What Is Generated

### 1. Model Structs And Schema Structs

For each model, Contractor generates:

- a base struct type mapping fields to Go types
- a schema struct type that embeds the base struct and adds the `K string` field mapped to `__k`

```go
type SignInMetadata struct {
	Browser   string `json:"browser"`
	UserAgent string `json:"userAgent"`
}

type SignInMetadataSchema struct {
	SignInMetadata
	K string `json:"__k"`
}
```

Why `K` (`__k`) matters:
- It identifies the model type in polymorphic arrays or JSON responses at runtime.

### 2. Generics

Generic models generate Go generic types using type parameters.

```go
type ApiResponse[T any] struct {
	Data T `json:"data"`
}

type ApiResponseSchema[T any] struct {
	ApiResponse[T]
	K string `json:"__k"`
}
```

### 3. REST Metadata

REST endpoints are emitted as variables that hold the path and method:

```go
var GetUser = struct {
	Path   string
	Method string
}{
	Path:   "/users/:id",
	Method: "GET",
}
```

### 4. Custom Errors

Errors declared in the contract are emitted as Go structs that implement the `error` interface. A constructor function is also provided.

```go
type InvalidEmailError struct {
	Message string
	Code    string
	Scope   string
	Status  int
}

func (e *InvalidEmailError) Error() string {
	return e.Message
}

func NewInvalidEmailError() *InvalidEmailError {
	return &InvalidEmailError{
		Message: "Invalid email address format",
		Code:    "INVALID_EMAIL",
	}
}
```

## Recommended Runtime Flow

Since Go performs strong type checking and Contractor emits strict types, you can directly unmarshal JSON into the generated structs.

```go
var response SignInMetadataSchema
err := json.Unmarshal(body, &response)
if err != nil {
	// handle error
}
```

### Note on Any Type
The Contractor type `Any` is mapped to Go's `interface{}` to ensure compatibility without relying on the `any` keyword.
