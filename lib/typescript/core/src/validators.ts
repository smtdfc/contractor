export interface IValidator<T> {
    validate(data: unknown, parent?: string): ValidateResult<T>
}

export interface ValidateResult<T> {
    success: boolean;
    data?: T;
    errors?: Record<string, string>;
}

export type ValidatorConstructor<T = unknown> = new () => IValidator<T>;

const validatorRegistry: Record<string, ValidatorConstructor> = {};

export function registerValidator<T>(modelHash: string, validatorClass: ValidatorConstructor<T>): ValidatorConstructor<T> {
    validatorRegistry[modelHash] = validatorClass;
    return validatorClass;
}

export function getRegisteredValidator(modelHash: string): ValidatorConstructor | undefined {
    return validatorRegistry[modelHash];
}

const globalScope = globalThis as typeof globalThis & {
    registerValidator?: typeof registerValidator;
    getRegisteredValidator?: typeof getRegisteredValidator;
};

if (typeof globalScope.registerValidator !== "function") {
    globalScope.registerValidator = registerValidator;
}

if (typeof globalScope.getRegisteredValidator !== "function") {
    globalScope.getRegisteredValidator = getRegisteredValidator;
}

export class Validator {
    private rules: Record<string, { type: string; errorMsg: string; config?: any }[]> = {};
    private currentKey: string = "";

    constructor() { }

    private ensureKey(name: string): void {
        if (!this.rules[name]) {
            this.rules[name] = [];
        }
    }

    private pushRule(type: string, errorMsg: string, config?: any): this {
        if (!this.currentKey) {
            return this;
        }

        this.ensureKey(this.currentKey);
        this.rules[this.currentKey]!.push({ type, errorMsg, config });
        return this;
    }

    private isEmpty(value: unknown): boolean {
        return value === undefined || value === null || value === "";
    }

    private toDate(value: unknown): Date | null {
        if (value instanceof Date) {
            return Number.isNaN(value.getTime()) ? null : value;
        }

        if (typeof value !== "string" && typeof value !== "number") {
            return null;
        }

        const date = new Date(value);
        return Number.isNaN(date.getTime()) ? null : date;
    }

    private getModelHash(value: unknown): string | undefined {
        if (!value || typeof value !== "object") {
            return undefined;
        }

        const record = value as Record<string, unknown>;
        const hashValue = record.__k ?? record._k;
        return typeof hashValue === "string" ? hashValue : undefined;
    }

    private beginField(name: string): this {
        this.currentKey = name;
        this.ensureKey(name);
        return this;
    }

    public key(name: string): this {
        return this.beginField(name);
    }

    public withGeneric(name: string, _value: unknown): this {
        this.beginField(name);
        this.rules[name]!.push({ type: "Generic", errorMsg: "" });
        return this;
    }

    public withArray(name: string, config: any): this {
        this.beginField(name);
        this.rules[name]!.push({ type: "Array", errorMsg: "", config });
        return this;
    }

    public withModel(name: string, validatorClass: any): this {
        this.beginField(name);
        this.rules[name]!.push({ type: "Nested", errorMsg: "", config: validatorClass });
        return this;
    }

    public Is(target: unknown, errorMsg: string): this {
        return this.pushRule("Is", errorMsg, target);
    }

    public Min(min: number, errorMsg: string): this {
        return this.pushRule("Min", errorMsg, min);
    }

    public Max(max: number, errorMsg: string): this {
        return this.pushRule("Max", errorMsg, max);
    }

    public Length(length: number, errorMsg: string): this {
        return this.pushRule("Length", errorMsg, length);
    }

    public MinLength(min: number, errorMsg: string): this {
        return this.pushRule("MinLength", errorMsg, min);
    }

    public MaxLength(max: number, errorMsg: string): this {
        return this.pushRule("MaxLength", errorMsg, max);
    }

    public Range(min: number, max: number, errorMsg: string): this {
        return this.pushRule("Range", errorMsg, { min, max });
    }

    public Matches(regex: RegExp | string, errorMsg: string): this {
        return this.pushRule("Matches", errorMsg, regex);
    }

    public Contains(substring: string, errorMsg: string): this {
        return this.pushRule("Contains", errorMsg, substring);
    }

    public StartsWith(substring: string, errorMsg: string): this {
        return this.pushRule("StartsWith", errorMsg, substring);
    }

    public EndsWith(substring: string, errorMsg: string): this {
        return this.pushRule("EndsWith", errorMsg, substring);
    }

    public In(values: unknown[], errorMsg: string): this {
        return this.pushRule("In", errorMsg, values);
    }

    public NotNull(errorMsg: string): this {
        return this.pushRule("NotNull", errorMsg);
    }

    public IsEmail(errorMsg: string): this {
        return this.pushRule("IsEmail", errorMsg);
    }

    public IsNumber(errorMsg: string): this {
        return this.pushRule("IsNumber", errorMsg);
    }

    public IsURL(errorMsg: string): this {
        return this.pushRule("IsURL", errorMsg);
    }

    public IsUUID(errorMsg: string): this {
        return this.pushRule("IsUUID", errorMsg);
    }

    public IsDate(errorMsg: string): this {
        return this.pushRule("IsDate", errorMsg);
    }

    public IsDateTime(errorMsg: string): this {
        return this.pushRule("IsDateTime", errorMsg);
    }

    public IsAlpha(errorMsg: string): this {
        return this.pushRule("IsAlpha", errorMsg);
    }

    public IsAlnum(errorMsg: string): this {
        return this.pushRule("IsAlnum", errorMsg);
    }

    public IsBool(errorMsg: string): this {
        return this.pushRule("IsBool", errorMsg);
    }

    public IsModel(errorMsg: string): this {
        return this.pushRule("IsModel", errorMsg);
    }

    public NestedValidate(errorMsg: string): this {
        const currentRules = this.rules[this.currentKey]!;
        const lastRule = currentRules[currentRules.length - 1];
        if (lastRule && (lastRule.type === "Nested" || lastRule.type === "Array")) {
            lastRule.errorMsg = errorMsg;
        }
        return this;
    }

    public execute<T>(data: any, parent = ""): ValidateResult<T> {
        let errors: Record<string, string> = {};
        const validatedData: any = { ...data };

        for (const key in this.rules) {
            const value = data ? data[key] : undefined;
            const keyRules = this.rules[key]!;
            const keyName = parent ? `${parent}.${key}` : key;
            const isValueEmpty = this.isEmpty(value);

            for (const rule of keyRules) {
                // Optional fields should not run type/format checks when empty.
                if (isValueEmpty) {
                    if (rule.type === "NotNull") {
                        errors[keyName] = rule.errorMsg;
                        break;
                    }

                    continue;
                }

                if (rule.type === "NotNull") {
                    continue;
                }

                if (rule.type === "Is" && value !== rule.config) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "Min" && (typeof value !== "number" || Number.isNaN(value) || value < rule.config)) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "Max" && (typeof value !== "number" || Number.isNaN(value) || value > rule.config)) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "Range") {
                    const range = rule.config as { min: number; max: number };
                    if (typeof value !== "number" || Number.isNaN(value) || value < range.min || value > range.max) {
                        errors[keyName] = rule.errorMsg;
                        break;
                    }
                }

                if (rule.type === "Length" && value?.length !== rule.config) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "MinLength" && (value?.length ?? -1) < rule.config) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "MaxLength" && (value?.length ?? Number.POSITIVE_INFINITY) > rule.config) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "Matches") {
                    const regex = typeof rule.config === "string" ? new RegExp(rule.config) : rule.config;
                    if (typeof value !== "string" || !regex.test(value)) {
                        errors[keyName] = rule.errorMsg;
                        break;
                    }
                }

                if (rule.type === "Contains" && (typeof value !== "string" || !value.includes(rule.config))) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "StartsWith" && (typeof value !== "string" || !value.startsWith(rule.config))) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "EndsWith" && (typeof value !== "string" || !value.endsWith(rule.config))) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "In" && (!Array.isArray(rule.config) || !rule.config.includes(value))) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "IsEmail") {
                    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
                    if (typeof value !== "string" || !emailRegex.test(value)) {
                        errors[keyName] = rule.errorMsg;
                        break;
                    }
                }

                if (rule.type === "IsNumber" && (typeof value !== "number" || Number.isNaN(value))) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "IsURL") {
                    try {
                        if (typeof value !== "string") {
                            throw new Error("invalid url");
                        }
                        new URL(value);
                    } catch {
                        errors[keyName] = rule.errorMsg;
                        break;
                    }
                }

                if (rule.type === "IsUUID" && (typeof value !== "string" || !/^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i.test(value))) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "IsDate" && !this.toDate(value)) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "IsDateTime" && !this.toDate(value)) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "IsAlpha" && (typeof value !== "string" || !/^[a-zA-Z]+$/.test(value))) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "IsAlnum" && (typeof value !== "string" || !/^[a-z0-9]+$/i.test(value))) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "IsBool" && typeof value !== "boolean") {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "IsModel" && (value === null || typeof value !== "object" || Array.isArray(value))) {
                    errors[keyName] = rule.errorMsg;
                    break;
                }

                if (rule.type === "Generic") {
                    const modelHash = this.getModelHash(value);
                    if (!modelHash) {
                        break;
                    }

                    const GenericValidator = getRegisteredValidator(modelHash);
                    if (!GenericValidator) {
                        errors[keyName] = rule.errorMsg || `VALIDATOR_NOT_REGISTERED:${modelHash}`;
                        break;
                    }

                    const subValidator = new GenericValidator();
                    const subResult = subValidator.validate(value, keyName);

                    if (!subResult.success) {
                        errors = { ...errors, ...subResult.errors };
                        break;
                    }
                }

                if (rule.type === "Nested" && value) {
                    const ValidatorClass = rule.config;
                    const subValidator = new ValidatorClass();
                    const subResult = subValidator.validate(value, keyName);

                    if (!subResult.success) {
                        errors = { ...errors, ...subResult.errors };
                        if (rule.errorMsg && !Object.keys(subResult.errors ?? {}).length) {
                            errors[keyName] = rule.errorMsg;
                        }
                        break;
                    }
                }

                if (rule.type === "Array" && value) {
                    if (!Array.isArray(value)) {
                        errors[keyName] = rule.errorMsg || "FIELD_MUST_BE_ARRAY";
                        break;
                    }

                    const itemConfig = rule.config;
                    let hasArrayError = false;

                    for (let i = 0; i < value.length; i++) {
                        const arrayKey = `${keyName}.${i}`;
                        const item = value[i];

                        if (typeof itemConfig === "function" || (itemConfig && itemConfig.prototype && itemConfig.prototype.validate)) {
                            const SubValidatorClass = itemConfig;
                            const subValidator = new SubValidatorClass();
                            const subResult = subValidator.validate(item, arrayKey);

                            if (!subResult.success) {
                                errors = { ...errors, ...subResult.errors };
                                hasArrayError = true;
                                break;
                            }
                        } else if (typeof itemConfig === "string") {
                            if (typeof item !== itemConfig) {
                                errors[arrayKey] = rule.errorMsg || `ITEM_MUST_BE_${itemConfig.toUpperCase()}`;
                                hasArrayError = true;
                                break;
                            }
                        }
                    }

                    if (hasArrayError) {
                        break;
                    }
                }
            }
        }

        const hasErrors = Object.keys(errors).length > 0;

        if (hasErrors) {
            return { success: false, errors };
        }

        return {
            success: true,
            data: validatedData as T
        };
    }
}