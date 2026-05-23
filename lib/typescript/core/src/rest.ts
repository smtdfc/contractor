import { IModelSchema } from "./model.js";

export type RequestBody = IModelSchema;
export type ResponseBody = IModelSchema;
export type RestMethod = "GET" | "POST" | "PUT";

export type RestMetadata<P extends string, M extends RestMethod, Req extends RequestBody, Res extends IModelSchema> = {
    path: P
    method: M
}

export type GetPath<T> = T extends RestMetadata<infer P, any, any, any> ? P : never;
export type GetMethod<T> = T extends RestMetadata<any, infer P, any, any> ? P : never;
export type GetRequest<T> = T extends RestMetadata<any, any, infer P, any> ? P : never;
export type GetResponse<T> = T extends RestMetadata<any, any, any, infer P> ? P : never;

