export type RestMethod = "GET" | "POST" | "PUT" | "DELETE" | "PATCH" | "OPTION";
export interface RestMetadata {
    path: string;
    method: RestMethod;
    queries: string[];
}

export type RestRequestBody<T> = T;
export type RestResponseBody<T> = T;