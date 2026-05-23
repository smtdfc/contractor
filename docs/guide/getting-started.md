# Getting Started

**Contractor** is a specialized Interface Definition Language (IDL) designed to enforce data integrity across distributed systems.

## Philosophy

In microservice architectures, the "contract" between services is often fragile. **Contractor** shifts the focus from manual implementation to **Schema-First development**. 

By using a single source of truth, it ensures that your data models, validation logic, and transformation rules are always in sync across your entire stack.

## Installation

You can install the Contractor CLI tool globally using Go:

```bash
go install github.com/smtdfc/contractor@latest
```
 

## Setup

Create a `contractor.json` in the root of your project:

```json
{
  "sourceDir": "./contracts",
  "extension": ".contract",
  "targets": [
    {
      "language": "typescript",
      "outDir": "./generated/ts"
    }
  ]
}
```

## Writing your first Contract

Create a directory named `contracts` (to match the `sourceDir` configured above) and add a file named `user.contract`:

```contractor
@CreateConstructor
model User {
    id: String
    username: String
    email: String
}

rest GetUser {
    path: "/api/v1/users"
    method: "GET"
    queries: ["id"]
}
```

This simple contract defines a `User` model and a basic REST endpoint.
The exact generated TypeScript artifacts (types, schema, mapper, validator, and runtime notes) are documented in `/guide/typescript-code-generation`.

## Generating Code

Once your `contractor.json` and `.contract` files are ready, you can run the code generator:

```bash
contractor generate --config contractor.json
```

If you only want to generate code for a specific target language defined in your config, you can use the `--lang` flag:

```bash
contractor generate --lang typescript
```

## TypeScript Runtime (`contractor-ts`)

If you are targeting TypeScript, the generated code relies on the `contractor-ts` runtime package for core utilities, base error classes, and validation logic.

Make sure to install it in your frontend or backend project:

```bash
npm install contractor-ts
# or
pnpm install contractor-ts
# or
yarn add contractor-ts
```
