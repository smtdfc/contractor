# Errors

Contractor provides a native way to define standard application errors, which helps in maintaining consistent error codes and HTTP statuses across distributed services.

## Defining an Error

Use the `error` keyword followed by the error name and a block specifying its attributes.

```contractor
error UserNotFoundError {
    message: "User not found"
    code: "USER_NOT_FOUND_ERR"
    status: 404
}
```

## Error Attributes

Errors support the following attributes:
- `message`: A human-readable description of the error (String).
- `code`: A unique, constant string identifier for the error (String).
- `status`: The HTTP status code associated with the error (Number).
- `scope`: (Optional) The scope or context of the error (String).

## Example

```contractor
error SignInInformationIncorrectError {
    message: "Sign-in Information Incorrect"
    code: "SIGNIN_INFO_INCORRECT_ERR"
    status: 400
}

error AuthorizationError {
    message: "Unauthenticated"
    code: "AUTHORIZATION_ERR"
    status: 401
}
```
