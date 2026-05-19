import { IContractModel, RestMetadata } from "contractor-ts";
import { HonoApp, validateMiddleware, ValidationTarget } from "./middleware.js";
import { Handler } from "hono";

export type RouteValidationOptions<T extends IContractModel> = {
    model: T,
    target: ValidationTarget
}


export function createRouteFromContract<T extends IContractModel>(app: HonoApp, rest: RestMetadata, validate?: RouteValidationOptions<T>, ...callback: Handler[]) {

    const methodName = rest.method.toLocaleLowerCase();
    if (validate) {
        app.on([methodName], [rest.path], validateMiddleware(validate.target, validate.model), ...callback);
        return
    }
    app.on([methodName], [rest.path], ...callback)
}