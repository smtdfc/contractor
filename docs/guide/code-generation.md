# Code Generation

For a complete and up-to-date TypeScript generation reference, including generated artifacts, runtime caveats, mapper and validator usage, and full examples, see:

- /guide/typescript/code-generation

Contractor's core utility lies in its ability to parse `.contract` files and generate idiomatic code for various target languages.

## The Generator Pipeline

1. **Lexer & Parser**: The IDL is parsed into an Abstract Syntax Tree (AST).
2. **Type Checker**: The AST is validated to ensure type correctness (e.g., ensuring generic types are provided, referenced models exist).
3. **IR Generation**: The AST is converted into an Intermediate Representation (IR), which simplifies the code generation process.
4. **Emitters**: Language-specific emitters traverse the IR and produce source code.

## TypeScript

The TypeScript emitter generates `.ts` files from your contracts.

Current TypeScript output includes:

- model types
- schema types with `__k`
- validator classes
- mapper classes
- generic mapper resolver types
- REST metadata constants

When you specify `typescript` as your target in `contractor.json`:

```json
{
  "targets": [
    {
      "language": "typescript",
      "outDir": "./generated/ts"
    }
  ]
}
```

Example contract:

```contractor
model SignInMetadata {
  browser: String
  @NotNull("FIELD_NOT_NULL")
  userAgent: String
}
```

Example generated TypeScript shape:

```typescript
export type SignInMetadata = {
  browser: string;
  userAgent: string;
};

export type SignInMetadataSchema = SignInMetadata & {
  __k: "8faa6476";
};

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

### Other Mappings

- **REST**: Maps to typed metadata constants.
- **Enums / Errors / Events**: Refer to current emitter implementation before relying on generation in your target version.

## Future Targets

Contractor's architecture is designed to support multiple languages. The code generation pipeline is designed to be extensible, allowing new emitters (e.g., for `Go`, `Java`, `Kotlin`, `C#`) to be added by implementing the `ProgramEmitter` interface.
