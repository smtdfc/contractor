# Validation

Contractor allows you to define validation rules directly on your model fields using annotations. These annotations are parsed by the IDL and can be used by the code generators to produce validation logic.

The TypeScript generator maps each supported annotation to a fluent method on `contractor-ts`'s `Validator` runtime. In generated code, each field becomes a chain such as `validator.key("email").NotNull("FIELD_NOT_NULL").IsEmail("INVALID_EMAIL")`.

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

## How It Works

1. The parser recognizes the annotation names listed above.
2. The generator extracts the annotation arguments into an intermediate representation.
3. The TypeScript emitter renders those rules into chained validator calls.
4. The `contractor-ts` runtime executes the rules and returns field-level errors.

Example of emitted TypeScript for one field:

```typescript
validator.key("email")
    .NotNull("FIELD_NOT_NULL")
    .IsEmail("INVALID_EMAIL");
```

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
