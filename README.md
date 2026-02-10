# Contractor: Type-Safe IDL & Code Generation Toolchain

**Contractor** is a specialized Interface Definition Language (IDL) designed to enforce data integrity across distributed systems. It provides a robust mechanism to define cross-platform data contracts, generating validated and idiomatic code for TypeScript and Go, eliminating the risks of manual synchronization.

### Core Philosophy

In microservice architectures, the "contract" between services is often fragile. **Contractor** shifts the focus from manual implementation to **Schema-First development**. By using a single source of truth, it ensures that your data models, validation logic, and transformation rules are always in sync across your entire stack.

### Key Technical Pillars

- **Strict Data Integrity**: Every generated model includes built-in validation logic to ensure incoming data adheres to the defined contract.
- **Predictable Mapping**: Explicitly handle field transformations (e.g., snake_case to camelCase) through `@Mapping` annotations, making API integration transparent.
- **Generic-First Design**: Native support for complex generic structures like `Response<T>`, ensuring type safety even in highly abstract data wrappers.
- **Developer Productivity**: Automates the generation of boilerplate code, including constructors, getters, setters, and mappers.

---
