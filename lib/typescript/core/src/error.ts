export type GeneratedErrorConstructor = new () => Error;
export type GeneratedErrorConstructorMap = Record<
    string,
    GeneratedErrorConstructor
>;


export abstract class ContractBaseError extends Error {
    static readonly TYPE = Symbol('ContractBaseError');
    readonly type = ContractBaseError.TYPE;

    constructor(public message: string) {
        super(message);
    }
}

export interface IContractError {
    status: number;
    message: string;
    code: string;
    type: symbol;
}

export function isContractError(err: unknown): err is IContractError {
    return typeof err === 'object' && err !== null && (err as any).type === ContractBaseError.TYPE;
}
