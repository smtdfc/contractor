import { IValidator, RestMethod, RestMetadata, IModelSchema } from "contractor-ts";
import { HonoApp, validateMiddleware, ValidationTarget } from "./middleware.js";
import { Handler } from "hono";

export type RouteValidationOptions<M, T extends IValidator<M>> = {
    validator: T,
    target: ValidationTarget
}


export function createRouteFromContract
    <
        M,
        T extends IValidator<M>,
        P extends string,
        K extends RestMethod,
        Req extends IModelSchema,
        Res extends IModelSchema
    >(app: HonoApp, rest: RestMetadata<P, K, Req, Res>, validate?: RouteValidationOptions<M, T>, ...callback: Handler[]) {

    const methodName = rest.method.toLocaleLowerCase();
    if (validate) {
        app.on([methodName], [rest.path], validateMiddleware(validate.target, validate.validator), ...callback);
        return
    }
    app.on([methodName], [rest.path], ...callback)
}