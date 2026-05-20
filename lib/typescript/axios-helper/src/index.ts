import axios, { isCancel, AxiosError, AxiosInstance } from "axios";
import { RestMetadata } from "contractor-ts"


type IResponseBody = {
    new(): any
    fromObject(): any
}

type IRequestBody = {
    new(): any
    fromObject(): any
}

export function createRequestFromContract<K, Q extends IResponseBody>(instance: AxiosInstance, rest: RestMetadata, data: K, config: any) {
    type AxiosMethods = 'get' | 'post' | 'put' | 'delete' | 'patch' | 'head' | 'options';
    const method = rest.method.toLowerCase() as AxiosMethods;
    const res = instance[method](rest.path, {}, {});
}

