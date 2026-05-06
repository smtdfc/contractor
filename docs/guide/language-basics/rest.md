# REST Endpoints

Contractor allows you to define RESTful API endpoints directly in your IDL, ensuring that your backend and frontend remain perfectly in sync regarding routing, methods, and payload structures.

## Defining an Endpoint

Use the `rest` keyword followed by an identifier, and a block specifying the route details.

```contractor
rest SignInWithEmail {
    path: "/api/v1/auth/sign-in/email"
    method: "POST"
}
```

## Endpoint Properties

REST endpoints support the following configuration options:

- `path`: The URL path of the endpoint (String literal).
- `method`: The HTTP method (e.g., `"GET"`, `"POST"`, `"PUT"`, `"DELETE"`).
- `queries`: An array of accepted query parameters (Array of strings).
- `request`: (Optional) The type of the request body payload.
- `response`: (Optional) The type of the response body payload.

## Example

```contractor
rest GetProfile {
    path: "/api/v1/auth/profile"
    method: "GET"
}

rest SearchUsers {
    path: "/api/v1/users"
    method: "GET"
    queries: ["query", "page", "limit"]
}
```
