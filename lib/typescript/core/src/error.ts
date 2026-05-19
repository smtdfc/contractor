export type GeneratedErrorConstructor = new () => Error;
export type GeneratedErrorConstructorMap = Record<
    string,
    GeneratedErrorConstructor
>;


export class ContractBaseError extends Error {
    constructor(msg: string) {
        super(msg);
    }
}