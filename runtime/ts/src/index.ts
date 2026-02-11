import { Validators } from "./validators.js";

export class ValidationError extends Error {
  public errors: Record<string, string>;

  constructor(errors: Record<string, string>) {
    super("ValidationError");
    this.errors = errors;
  }
}

export const ContractorRuntime = {
  Validators,
  ValidationError,
};
