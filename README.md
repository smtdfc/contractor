# Contractor: High-Performance IDL & Code Generation Toolchain

**Contractor** is a specialized Interface Definition Language (IDL) and runtime ecosystem designed to ensure strict type safety and data integrity across distributed systems. It allows developers to define data contracts once and generate idiomatic, validated code for multiple target environments, primarily focusing on TypeScript and Go.

## 1. Overview

In modern microservices and client-server architectures, maintaining data consistency is a recurring challenge. **Contractor** functions as the "Single Source of Truth." By defining models in a platform-agnostic IDL, you eliminate manual boilerplate and synchronization issues between backend and frontend teams.

### Key Features

- **Declarative Schema**: Simple, human-readable syntax for defining complex data structures.
- **Annotation-Driven Logic**: Use annotations like `@Data`, `@Optional`, or `@IsEmail` to inject behaviors and validation rules automatically.
- **Type Safety & Generics**: Native support for Generic types (e.g., `Response<T>`) and primitive mapping between languages.
- **Built-in Validation**: Generates static `validate()` methods that integrate seamlessly with the `ContractorRuntime` library.

---

## 2. Specification Syntax

Define your models in `.ctr` files using the Contractor IDL. The syntax is designed to be familiar to users of Java, C#, or TypeScript.

```java
@CreateConstructor
@Data
model LoginDTO {
    @Private
    @Optional
    @IsEmail("Invalid email format")
    String email

    @Required
    String password
}

@Data
model Response<T> {
    Array<T> data
    T message
}

```

---

## 3. Generated Code Architecture (TypeScript)

The `TypescriptGenerator` processes the AST (Abstract Syntax Tree) to produce production-ready classes. Below is the generated output for the `LoginDTO` model:

```typescript
import { ContractorRuntime } from "contractor";

export class LoginDTO {
  private email?: string;
  public password: string;

  constructor(email: string, password: string) {
    this.email = email;
    this.password = password;
  }

  public getEmail(): string {
    return this.email;
  }

  public setEmail(v: string): void {
    this.email = v;
  }

  public getPassword(): string {
    return this.password;
  }

  public setPassword(v: string): void {
    this.password = v;
  }

  public static validate(obj: any): ContractorRuntime.ValidationError | null {
    const errors: Record<string, string> = {};

    if (!ContractorRuntime.Validators.IsRequired(obj.password)) {
      errors["password"] = "password is required";
    }

    if (Object.keys(errors).length > 0) {
      return new ContractorRuntime.ValidationError(errors);
    }
    return null;
  }
}

export class Response<T> {
  public Data: Array<T>;
  public Hello: T;

  public getData(): Array<T> {
    return this.Data;
  }

  public setData(v: Array<T>): void {
    this.Data = v;
  }

  public getHello(): T {
    return this.Hello;
  }

  public setHello(v: T): void {
    this.Hello = v;
  }

  public static validate(obj: any): ContractorRuntime.ValidationError | null {
    const errors: Record<string, string> = {};

    if (!ContractorRuntime.Validators.IsRequired(obj.Data)) {
      errors["Data"] = "Data is required";
    }
    if (!ContractorRuntime.Validators.IsRequired(obj.Hello)) {
      errors["Hello"] = "Hello is required";
    }

    if (Object.keys(errors).length > 0) {
      return new ContractorRuntime.ValidationError(errors);
    }
    return null;
  }
}
```

---

## 4. Technical Implementation

The toolchain is implemented in **Go** for high-performance parsing and code synthesis.

1. **Lexer/Parser**: Tokenizes the IDL source and constructs an Abstract Syntax Tree (AST).
2. **Generator Engine**: Utilizes a specialized `CodeBuffer` to manage indentation levels and stream-based code synthesis, ensuring clean and readable output.
3. **Refinement**: Automatically handles casing conversions (e.g., converting IDL field names to PascalCase for TypeScript getters).

---
