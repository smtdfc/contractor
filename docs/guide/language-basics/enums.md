# Enums

Enums define a fixed set of named values.

## Basic Syntax

Use `enum` followed by the enum name and members.
Members can be separated by commas and/or newlines.

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

## Common Mistakes

1. Duplicate member names in the same enum.
2. Empty enum declarations.
3. Using enum type before it is declared incorrectly in malformed files.
