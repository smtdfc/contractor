# Models

Models are the core building blocks of a Contractor schema.
They define data shape, field types, optionality, and validation annotations.

## Basic Syntax

Use the `model` keyword followed by the model name and a block containing fields.

```contractor
model User {
    id: String
    username: String
    email: String
}
```

Field format:

```contractor
fieldName: Type
optionalField?: Type
```

## Optional Fields

You can mark fields as optional by appending `?` to the field name.

```contractor
model User {
    id: String
    email: String
    avatar?: String
    displayName?: String
}
```

Notes:

- Optional means the field may be missing.
- Runtime validator skips non-`NotNull` rules when an optional field is empty.

## Generics

Models support generic type parameters.

```contractor
model ApiResponse<T> {
    data: T
}

model ProfileResponse {
    payload: ApiResponse<String>
}
```

Generic type arguments are validated by the type checker. If a type expects generic parameters, you must provide all required arguments.

## Built-in Types

Contractor supports the following built-in types:

- `String`
- `Int`
- `Float`
- `Bool`
- `Array<T>`
- `Null`
- `Any`

Examples:

```contractor
model Metrics {
    values: Array<Float>
    tags: Array<String>
    extra: Any
}
```

## Annotations

Models and fields can be decorated with annotations.
Model-level annotations configure generation behavior, while field-level annotations control validation rules.

For the full validation annotation list and argument signatures, see `/guide/validation`.

### `@CreateConstructor`

Marks model intent for constructor-related generation behavior.

```contractor
@CreateConstructor
model Token {
    accessToken: String
}
```

### `@Mapper`

Marks model intent for mapper-related generation behavior.

```contractor
@Mapper
model SignInMetadata {
    browser: String
    userAgent: String
}
```

Important:

- Current type checker recognizes `@Mapper`.
- `@CreateMapper` is not a built-in annotation in the current parser.

## Full Example

```contractor
@CreateConstructor
@Mapper
model SignInRequest {
    @NotNull("FIELD_NOT_NULL")
    @IsEmail("INVALID_EMAIL")
    email: String

    @NotNull("FIELD_NOT_NULL")
    password: String

    rememberMe?: Bool
}
```

## Common Mistakes

1. Using an unknown annotation name.
2. Missing generic type arguments (for example `Array` instead of `Array<String>`).
3. Applying `@NestedValidate` on a non-model field.
