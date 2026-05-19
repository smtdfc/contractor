# Contractor TypeScript Core

Type-safe runtime helpers and shared types for Contractor-generated TypeScript contracts.

## Overview

`contractor-ts` provides a comprehensive `Validator` utility for validating data against the contracts defined in your Contractor IDL. These validators are used by generated code to ensure data integrity and provide meaningful error messages when validation fails.

## Installation

```bash
npm install contractor-ts
# or
pnpm add contractor-ts
```

## Usage

The main export is the `Validator` object, which contains a collection of validation functions. Each validator follows a consistent pattern:

```typescript
import { Validator } from "contractor-ts";

// Returns null if validation passes, or the error message if it fails
const result = Validator.IsEmail("test@example.com", "Invalid email");
// result === null (validation passed)

const result = Validator.IsEmail("invalid-email", "Invalid email");
// result === 'Invalid email' (validation failed)
```

## Available Validators

### Basic Comparisons

- **`Is(value, target, errorMsg)`** - Check if value equals target
- **`In(value, list, errorMsg)`** - Check if value is in a list

### Number Validators

- **`Min(value, min, errorMsg)`** - Check if number is >= min
- **`Max(value, max, errorMsg)`** - Check if number is <= max
- **`Range(value, [min, max], errorMsg)`** - Check if number is within range
- **`IsNumber(value, errorMsg)`** - Check if value is a number

### String Length Validators

- **`Length(value, len, errorMsg)`** - Check if length equals len
- **`MinLength(value, min, errorMsg)`** - Check if length >= min
- **`MaxLength(value, max, errorMsg)`** - Check if length <= max

### String Pattern Validators

- **`Matches(value, regex, errorMsg)`** - Check if string matches regex
- **`Contains(value, sub, errorMsg)`** - Check if string contains substring
- **`StartsWith(value, sub, errorMsg)`** - Check if string starts with substring
- **`EndsWith(value, sub, errorMsg)`** - Check if string ends with substring

### Format Validators

- **`IsEmail(value, errorMsg)`** - Validate email format
- **`IsURL(value, errorMsg)`** - Validate URL format
- **`IsUUID(value, errorMsg)`** - Validate UUID format
- **`IsDate(value, errorMsg)`** - Validate date string
- **`IsDateTime(value, errorMsg)`** - Validate datetime string
- **`IsAlpha(value, errorMsg)`** - Check if string contains only letters
- **`IsAlnum(value, errorMsg)`** - Check if string contains only alphanumeric characters

### Type Validators

- **`IsBool(value, errorMsg)`** - Check if value is boolean
- **`NotNull(value, errorMsg)`** - Check if value is not null/undefined
- **`IsModel(value, errorMsg)`** - Check if value is an object

### Nested Validation

- **`NestedValidate(value, errorMsg)`** - Validate nested objects that have a `validate()` method

## Example

```typescript
import { Validator } from "contractor-ts";

// Validate multiple fields
const errors: Record<string, string> = {};

const email = "user@example.com";
const error = Validator.IsEmail(email, "Email is invalid");
if (error) errors.email = error;

const age = 25;
const ageError = Validator.Range(
  age,
  [0, 120],
  "Age must be between 0 and 120",
);
if (ageError) errors.age = ageError;
```

## License

MIT

## Repository

[https://github.com/smtdfc/contractor](https://github.com/smtdfc/contractor)
