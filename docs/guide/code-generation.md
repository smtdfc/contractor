# Code Generation

Contractor's core utility lies in its ability to parse `.contract` files and generate idiomatic code for various target languages.

## The Generator Pipeline

1. **Lexer & Parser**: The IDL is parsed into an Abstract Syntax Tree (AST).
2. **Type Checker**: The AST is validated to ensure type correctness (e.g., ensuring generic types are provided, referenced models exist).
3. **IR Generation**: The AST is converted into an Intermediate Representation (IR), which simplifies the code generation process.
4. **Emitters**: Language-specific emitters traverse the IR and produce source code.

## TypeScript

The TypeScript emitter generates `.ts` files based on your models and definitions.

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

By default, a `model` in Contractor maps to a standard TypeScript `class`.

**Contractor:**
```contractor
model User {
    id: String
    username: String
    email: String
}
```

**Generated TypeScript:**
```typescript
export class User {
    id: string;
    username: string;
    email: string;
}
```

If you use the `@CreateConstructor` annotation, the emitter will generate a constructor method within the class, integrating validation logic for the fields based on their annotations.

**Contractor:**
```contractor
@CreateConstructor
model User {
    id: String
    @IsEmail("INVALID_EMAIL")
    email: String
}
```

**Generated TypeScript (Conceptual):**
```typescript
import { Validator } from "contractor-ts";

export class User {
    id: string;
    email: string;

    constructor(data: any) {
        // Validation logic using contractor-ts
        this.id = data.id;
        this.email = data.email; // Validated for @IsEmail
    }
}
```

*(Note: The exact output structure depends on the internal implementation of the emitter, but it consistently maps your data constraints into executable TypeScript code).*

### Other Mappings

- **Enums**: Map directly to TypeScript `enum` declarations.
- **Errors**: Map to custom error classes extending `Error` (or a base error from `contractor-ts`).
- **REST / Events**: Map to strongly-typed constants or function signatures that you can use with your HTTP clients or event buses.

## Future Targets

Contractor's architecture is designed to support multiple languages. The code generation pipeline is designed to be extensible, allowing new emitters (e.g., for `Go`, `Java`, `Kotlin`, `C#`) to be added by implementing the `ProgramEmitter` interface.
