# Enums

Enums (Enumerations) allow you to define a type by enumerating its possible values.

## Defining an Enum

Use the `enum` keyword followed by the enum name and a block containing its members. Members are separated by commas.

```contractor
enum Role {
    User, Admin, SuperAdmin
}
```

## Usage

Once defined, enums can be used as types for model fields.

```contractor
model User {
    id: String
    username: String
    role: Role
}
```
