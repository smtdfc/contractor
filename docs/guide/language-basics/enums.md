# Enums

Enums define a fixed set of named values.

## Basic Syntax

Use `enum` followed by the enum name and a block of members.
Members are identifiers, not string literals.
Commas are optional between members, and newlines are allowed.

```contractor
enum Role {
    User,
    Admin,
    SuperAdmin
}
```

## Usage In Models

Once declared, enums can be used as field types.

```contractor
model User {
    id: String
    username: String
    role: Role
}
```

## Rules Checked By Type Checker

1. Enum name is required.
2. Enum must have at least one member.
3. Member names must be unique.

The parser also rejects non-identifier members.

## Full Example

```contractor
enum Permission {
    Read,
    Write,
    Delete
}

model Policy {
    id: String
    defaultPermission: Permission
    allowed: Array<Permission>
}
```

## Invalid Examples

```contractor
enum BadRole {
    User,
    User
}
```

This fails because the same member is declared twice.

```contractor
enum Empty {
}
```

This fails because an enum must have at least one member.

```contractor
enum WrongMember {
    "Admin"
}
```

This fails because enum members must be identifiers.

## Common Mistakes

1. Duplicate member names in the same enum.
2. Empty enum declarations.
3. Writing members as strings instead of identifiers.
