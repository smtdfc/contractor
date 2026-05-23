# REST Endpoints

REST declarations define endpoint metadata in IDL so service/client contracts stay consistent.

## Basic Syntax

Use `rest` followed by an identifier and a property block.

```contractor
rest SignInWithEmail {
    path: "/api/v1/auth/sign-in/email"
    method: "POST"
}
```

## Properties

Supported properties:

- `path` (required): string literal
- `method` (required): string literal, one of `GET`, `POST`, `PUT`, `PATCH`, `DELETE` (case-insensitive)
- `queries` (optional): array literal of strings
- `requestBody` (optional): user-defined type or `Null`
- `responseBody` (optional): user-defined type or `Null`

Type checker notes:

1. `method` and `path` are required.
2. `requestBody` and `responseBody` do not accept primitive built-in types.
3. `GET` with non-null `requestBody` is allowed but emits a warning.
4. Unknown or duplicate properties are rejected.

## Example

```contractor
rest GetProfile {
    path: "/api/v1/auth/profile"
    method: "GET"
    queries: ["includeRoles"]
    responseBody: ProfileResponse
}

rest SignInWithEmail {
    path: "/api/v1/auth/sign-in/email"
    method: "POST"
    requestBody: SignInWithEmailRequest
    responseBody: SignInResponse
}
```

## Common Mistakes

1. Using `request`/`response` instead of `requestBody`/`responseBody`.
2. Using non-string values inside `queries`.
3. Supplying built-in primitive type directly to `requestBody` or `responseBody`.
