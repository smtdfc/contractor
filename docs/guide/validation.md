# Validation

Contractor allows you to define validation rules directly on your model fields using annotations. These annotations are parsed by the IDL and can be used by the code generators to produce validation logic.

## Supported Validators

The following validation annotations are natively recognized by the Contractor parser:

### General
- `@Is(type: String)`
- `@NotNull(message?: String)`
- `@In(values: Array)`

### Strings
- `@Length(min: Number, max: Number)`
- `@MinLength(min: Number)`
- `@MaxLength(max: Number)`
- `@Matches(regex: String)`
- `@Contains(substring: String)`
- `@StartsWith(substring: String)`
- `@EndsWith(substring: String)`
- `@IsEmail(message?: String)`
- `@IsURL()`
- `@IsUUID()`
- `@IsAlpha()`
- `@IsAlnum()`

### Numbers
- `@Min(value: Number)`
- `@Max(value: Number)`
- `@Range(min: Number, max: Number)`
- `@IsNumber()`

### Dates
- `@IsDate()`
- `@IsDateTime()`

### Booleans
- `@IsBool()`

### Relational / Nested
- `@IsModel()`
- `@NestedValidate(message?: String)`: Instructs the validator to deeply validate a nested model.

## Usage Example

Validators take arguments which can be used to customize error messages or validation thresholds.

```contractor
@CreateConstructor
model SignInWithEmailRequest {
    @NotNull("FIELD_NOT_NULL")
    @IsEmail("INVALID_EMAIL")
    email: String

    @NotNull("FIELD_NOT_NULL")
    @MinLength(8)
    password: String

    @NotNull("FIELD_NOT_NULL")
    @NestedValidate("INVALID_DATA")
    metadata: SignInMetadata
}
```
