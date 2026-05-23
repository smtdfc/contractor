# Events

Events describe asynchronous messages for pub/sub or event-driven flows.

## Basic Syntax

Use `event` followed by the event identifier and a property block.

```contractor
event SignUpEvent {
    name: "auth.sign-up"
    payload: UserCreatedPayload
}
```

## Properties

Events support:

- `name` (required): string literal
- `payload` (required): type reference

Type checker rules:

1. `name` is required and must be a string literal.
2. `payload` is required and must be a valid type.
3. Unknown event properties are rejected.

## Full Example

```contractor
model UserCreatedPayload {
    id: String
    email: String
}

event UserCreated {
    name: "user.created"
    payload: UserCreatedPayload
}
```

## Common Mistakes

1. Treating `name` as optional.
2. Using non-string value for `name`.
3. Missing `payload`.
