export type RestMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH" | "OPTION";

interface IClassConstructor {
    new(...args: any[]): any;
    fromObject(obj: any): any;
}


export interface RestMetadata<P extends string, M extends RestMethod, Req extends IClassConstructor, Res extends IClassConstructor> { // path, method, request , response
    path: string;
    method: RestMethod;
    queries: string[];
    requestBody: Req;
    responseBody: Res;
}

export type RestRequestBody<T> = T;
export type RestResponseBody<T> = T;
