export type GeneratedErrorConstructor = new () => Error;
export type GeneratedErrorConstructorMap = Record<
    string,
    GeneratedErrorConstructor
>;
