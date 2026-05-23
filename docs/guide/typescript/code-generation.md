# TypeScript Code Generation

This guide explains exactly what Contractor currently generates for TypeScript, how to use generated classes safely, and what to watch out for in production code.

## Scope

Current TypeScript output focuses on:

- model types
- schema types with discriminator `__k`
- one validator class per model
- one mapper class per model
- generic resolver types for generic models
- REST metadata constants

## Generation Pipeline

1. Parse `.contract` files into AST.
2. Type-check models, references, generics, and annotations.
3. Convert AST to IR.
4. Emit TypeScript source from IR.

## What Is Generated

### 1. Model Type And Schema Type

For each model, Contractor generates:

- a plain model type
- a schema type that includes `__k`

```ts
export type SignInMetadata = {
  browser: string;
  userAgent: string;
};

export type SignInMetadataSchema = SignInMetadata & {
  __k: "8faa6476";
};
```

Why `__k` matters:

- it identifies model type at runtime
- generic and nested validation can use it to resolve validators

### 2. Validator Class

Each model gets a validator class and is auto-registered in runtime registry.

```ts
export class SignInMetadataValidator implements IValidator {
  static readonly modelHash = "8faa6476";

  validate(data: any, parent = ""): ValidateResult<SignInMetadataSchema> {
    const validator = new Validator();
    validator.key("userAgent").NotNull("FIELD_NOT_NULL");
    return validator.execute(data, parent);
  }
}

registerValidator("8faa6476", SignInMetadataValidator);
```

### 3. Mapper Class

Each model gets a mapper that normalizes input and injects `__k`.

```ts
export class SignInMetadataMapper {
  static from(data: any): SignInMetadataSchema {
    const source = (data ?? {}) as Record<string, any>;

    return {
      __k: "8faa6476",
      browser: source.browser,
      userAgent: source.userAgent,
    };
  }
}
```

Important:

- mapper transforms shape and sets discriminator
- validator enforces constraints
- mapper does not replace validator

### 4. Generic Resolver Types

Generic models generate resolver types and require resolver functions when mapping.

```ts
export type ApiResponseGenericMappers<T> = {
  T: (value: any, path?: string) => T;
};

export class ApiResponseMapper {
  static from<T>(
    data: any,
    genericMappers: ApiResponseGenericMappers<T>,
  ): ApiResponseSchema<T> {
    const source = (data ?? {}) as Record<string, any>;

    return {
      __k: "7a862810",
      data: genericMappers.T(source.data, "data"),
    };
  }
}
```

Why resolver is required for generic mapping:

- generic type parameters exist only at compile-time
- runtime cannot infer what `T` is from `any` in a reliable way

### 5. REST Metadata

REST endpoints are emitted as typed metadata constants used by clients/helpers.

## Recommended Runtime Flow

Always use this order:

1. map raw input to schema
2. validate mapped schema
3. use `result.data` only when validation succeeds

Example helper:

```ts
function mapAndValidate<TSchema>(
  mapper: () => TSchema,
  validator: {
    validate: (data: any, parent?: string) => ValidateResult<TSchema>;
  },
): ValidateResult<TSchema> {
  const mapped = mapper();
  return validator.validate(mapped);
}
```

## Usage Examples

### Non-Generic Model

```ts
const raw = {
  browser: "Chrome",
  userAgent: "Mozilla/5.0",
};

const mapped = SignInMetadataMapper.from(raw);
const result = new SignInMetadataValidator().validate(mapped);

if (!result.success) {
  console.error(result.errors);
}
```

### Nested Model

```ts
const raw = {
  email: "user@example.com",
  password: "secret",
  metadata: {
    browser: "Chrome",
    userAgent: "Mozilla/5.0",
  },
};

const mapped = SignInWithEmailRequestMapper.from(raw);
const result = new SignInWithEmailRequestValidator().validate(mapped);

if (!result.success) {
  console.error(result.errors);
}
```

### Generic Model With Primitive Payload

```ts
const raw = { data: "123" };

const mapped = ApiResponseMapper.from<string>(raw, {
  T: (value) => String(value),
});

const result = new ApiResponseValidator().validate(mapped);
```

### Generic Model With Model Payload

```ts
const raw = {
  data: {
    browser: "Chrome",
    userAgent: "Mozilla/5.0",
  },
};

const mapped = ApiResponseMapper.from<SignInMetadataSchema>(raw, {
  T: (value) => SignInMetadataMapper.from(value),
});

const result = new ApiResponseValidator().validate(mapped);
```

## Full End-To-End Example

### Contract

```contractor
@CreateConstructor
model SignInMetadata {
  browser: String
  @NotNull("FIELD_NOT_NULL")
  userAgent: String
}

@CreateConstructor
model SignInWithEmailRequest {
  @NotNull("FIELD_NOT_NULL")
  @IsEmail("INVALID_EMAIL")
  email: String

  @NotNull("FIELD_NOT_NULL")
  password: String

  @NotNull("FIELD_NOT_NULL")
  @NestedValidate("INVALID_DATA")
  metadata: SignInMetadata
}

@CreateConstructor
model ApiResponse<T> {
  data: T
}
```

### Application Code

```ts
import {
  SignInMetadataMapper,
  SignInWithEmailRequestMapper,
  SignInWithEmailRequestValidator,
  ApiResponseMapper,
  ApiResponseValidator,
  type SignInMetadataSchema,
} from "./shared/contracts/src/wind-work-auth/index";

const signInRaw = {
  email: "user@example.com",
  password: "secret",
  metadata: {
    browser: "Chrome",
    userAgent: "Mozilla/5.0",
  },
};

const signInMapped = SignInWithEmailRequestMapper.from(signInRaw);
const signInValidated = new SignInWithEmailRequestValidator().validate(
  signInMapped,
);

if (!signInValidated.success) {
  console.error("SignIn invalid", signInValidated.errors);
}

const responseRaw = {
  data: {
    browser: "Chrome",
    userAgent: "Mozilla/5.0",
  },
};

const responseMapped = ApiResponseMapper.from<SignInMetadataSchema>(
  responseRaw,
  {
    T: (value) => SignInMetadataMapper.from(value),
  },
);

const responseValidated = new ApiResponseValidator().validate(responseMapped);

if (!responseValidated.success) {
  console.error("ApiResponse invalid", responseValidated.errors);
}
```

## Notes And Caveats

1. `__k` is required for discriminator-based model resolution.
2. If you skip mapper and build objects manually, make sure model objects still include `__k`.
3. Optional fields now skip non-`NotNull` rules when value is empty (`undefined`, `null`, empty string).
4. For generic values:
   - object values with `__k` can be deep-validated via registry
   - values without `__k` are not deep-validated as model objects
5. Mapper generation is currently available by default in emitted model output.

## Practical Checklist

1. Map first, validate second.
2. Reuse mapper output for nested and generic values.
3. Provide explicit generic resolvers for all type parameters.
4. Add tests for primitive generic, model generic, and invalid nested payload cases.
