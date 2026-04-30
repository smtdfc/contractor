export const Validator = {
  Is: (value: unknown, target: unknown, errorMsg: string) => {
    return value === target ? null : errorMsg;
  },

  Min: (value: unknown, min: number, errorMsg: string) => {
    return typeof value === "number" && value >= min ? null : errorMsg;
  },
  Max: (value: unknown, max: number, errorMsg: string) => {
    return typeof value === "number" && value <= max ? null : errorMsg;
  },
  Range: (value: unknown, [min, max]: [number, number], errorMsg: string) => {
    return typeof value === "number" && value >= min && value <= max
      ? null
      : errorMsg;
  },

  Length: (value: any, len: number, errorMsg: string) => {
    return value?.length === len ? null : errorMsg;
  },
  MinLength: (value: any, min: number, errorMsg: string) => {
    return value?.length >= min ? null : errorMsg;
  },
  MaxLength: (value: any, max: number, errorMsg: string) => {
    return value?.length <= max ? null : errorMsg;
  },

  Matches: (value: unknown, regex: RegExp, errorMsg: string) => {
    return typeof value === "string" && regex.test(value) ? null : errorMsg;
  },
  Contains: (value: string, sub: string, errorMsg: string) => {
    return value?.includes(sub) ? null : errorMsg;
  },
  StartsWith: (value: string, sub: string, errorMsg: string) => {
    return value?.startsWith(sub) ? null : errorMsg;
  },
  EndsWith: (value: string, sub: string, errorMsg: string) => {
    return value?.endsWith(sub) ? null : errorMsg;
  },

  In: (value: unknown, list: unknown[], errorMsg: string) => {
    return list.includes(value) ? null : errorMsg;
  },

  IsEmail: (value: string, errorMsg: string) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(value) ? null : errorMsg;
  },
  IsNumber: (value: unknown, errorMsg: string) => {
    return typeof value === "number" && !isNaN(value) ? null : errorMsg;
  },
  IsURL: (value: string, errorMsg: string) => {
    try {
      new URL(value);
      return null;
    } catch {
      return errorMsg;
    }
  },
  IsUUID: (value: string, errorMsg: string) => {
    const uuidRegex =
      /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
    return uuidRegex.test(value) ? null : errorMsg;
  },
  IsDate: (value: string, errorMsg: string) => {
    return !isNaN(Date.parse(value)) ? null : errorMsg;
  },
  IsDateTime: (value: string, errorMsg: string) => {
    return !isNaN(Date.parse(value)) ? null : errorMsg;
  },
  IsAlpha: (value: string, errorMsg: string) => {
    return /^[a-zA-Z]+$/.test(value) ? null : errorMsg;
  },
  IsAlnum: (value: string, errorMsg: string) => {
    return /^[a-z0-9]+$/i.test(value) ? null : errorMsg;
  },

  NotNull: (value: unknown, errorMsg: string) => {
    return value !== null && value !== undefined ? null : errorMsg;
  },
  IsBool: (value: unknown, errorMsg: string) => {
    return typeof value === "boolean" ? null : errorMsg;
  },

  IsModel: (value: unknown, errorMsg: string) => {
    return value !== null && typeof value === "object" ? null : errorMsg;
  },
  NestedValidate: (value: unknown, errorMsg: string) => {
    if (!value || typeof value !== "object") {
      return errorMsg;
    }

    const validate = (value as any).constructor?.validate;
    if (typeof validate !== "function") {
      return errorMsg;
    }

    const details = validate(value);
    return details && Object.keys(details).length > 0 ? errorMsg : null;
  },
};

export type RestMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH" | "OPTION";
export interface RestMetadata {
  path: string;
  method: RestMethod;
  queries: string[];
}

export type RestRequestBody<T> = T;
export type RestResponseBody<T> = T;
export type GeneratedErrorConstructor = new () => Error;
export type GeneratedErrorConstructorMap = Record<
  string,
  GeneratedErrorConstructor
>;
export type GeneratedValidationDetails = Record<string, string[]>;
export interface EventMetadata {
  name: string;
  method: string;
}

export type EventPayload<T> = T;
