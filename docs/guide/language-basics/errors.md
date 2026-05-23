# Errors

Errors let you define consistent domain/application errors in one place.
They are useful for keeping error code and status conventions aligned across services.

## Basic Syntax

Use `error` followed by the error name and a property block.

```contractor
error UserNotFoundError {
    message: "User not found"
    code: "USER_NOT_FOUND_ERR"
    status: 404
}
```

## Properties

Supported properties:

- `message` (required): string literal
- `code` (optional): string literal
- `scope` (optional): string literal
- `status` (optional): string literal or number literal

Type checker rules:

1. `message` is required.
2. `code` must be string if provided.
3. `scope` must be string if provided.
4. `status` must be string or number if provided.

The parser only accepts literal values here. You cannot use a model field, variable, or expression.

## Example

```contractor
error SignInInformationIncorrectError {
    message: "Sign-in information incorrect"
    code: "SIGNIN_INFO_INCORRECT_ERR"
    status: 400
}

error AuthorizationError {
    message: "Unauthenticated"
    code: "AUTHORIZATION_ERR"
    status: "401"
    scope: "auth"
}
```

## Invalid Examples

```contractor
error BrokenError {
    code: "BROKEN_ERR"
    status: 500
}
```

This fails because `message` is required.

```contractor
error BadStatusError {
    message: "Bad"
    status: true
}
```

This fails because `status` must be a string or number literal.

```contractor
error BadCodeError {
    message: "Bad"
    code: 123
}
```

This fails because `code` must be a string literal.

## Common Mistakes

1. Omitting `message`.
2. Using non-literal values for properties.
3. Using unknown properties in error blocks.
