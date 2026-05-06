# Models

Models are the core building blocks of a Contractor IDL schema. They define the shape of your data, including fields, types, and constraints.

## Defining a Model

Use the `model` keyword followed by the model name and a block containing fields.

```contractor
model User {
    id: String
    username: String
    email: String
}
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

## Generics

Models support generic type parameters.

```contractor
model User<T> {
    id: T
    username: String
}

model GetProfileResponse {
    user: User<String>
}
```

## Built-in Types

Contractor supports the following built-in types:
- `String`
- `Int`
- `Float`
- `Bool`
- `Array`
- `Null`
- `Any`

## Annotations

Models can be decorated with annotations to configure code generation behavior.

### `@CreateConstructor`

Generates a constructor function or class constructor for the model.

```contractor
@CreateConstructor
model Token {
    accessToken: String
}
```

### `@CreateMapper`

Instructs the code generator to emit data mapping utilities for this model.
