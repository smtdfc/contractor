export const Validators = {
  IsRequired(value: any): boolean {
    return value !== null && value !== undefined;
  },
  IsEmail(value: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(value);
  },
  IsNumber(value: any): boolean {
    return typeof value === "number" && !isNaN(value);
  },
  IsInt(value: any): boolean {
    return Number.isInteger(value);
  },
  IsFloat(value: any): boolean {
    return typeof value === "number" && !Number.isInteger(value);
  },
  IsBoolean(value: any): boolean {
    return typeof value === "boolean";
  },
  IsString(value: any): boolean {
    return typeof value === "string";
  },
  IsArray(value: any): boolean {
    return Array.isArray(value);
  },
  IsUrl(value: string): boolean {
    try {
      new URL(value);
      return true;
    } catch {
      return false;
    }
  },
  IsUUID(value: string): boolean {
    const uuidRegex =
      /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
    return uuidRegex.test(value);
  },
  IsNotEmpty(value: string | any[]): boolean {
    if (typeof value === "string") {
      return value.trim().length > 0;
    }
    if (Array.isArray(value)) {
      return value.length > 0;
    }
    return false;
  },
  Max(value: number, max: number): boolean {
    return value <= max;
  },
  Min(value: number, min: number): boolean {
    return value >= min;
  },
  Length(value: string, length: number): boolean {
    return value.length === length;
  },
  MinLength(value: string, minLength: number): boolean {
    return value.length >= minLength;
  },
  MaxLength(value: string, maxLength: number): boolean {
    return value.length <= maxLength;
  },
  ArrayMinLength(value: any[], minLength: number): boolean {
    return value.length >= minLength;
  },
  ArrayMaxLength(value: any[], maxLength: number): boolean {
    return value.length <= maxLength;
  },
};
