import { IValidator } from "contractor-ts";
import { Hono, Context, Next, Handler } from "hono"
import { ValidationError } from "./errors.js";

export type HonoApp = Hono<{ Bindings: any, Variables: any }>;
export type AppContext = Context<{ Bindings: any, Variables: any }>;


export type ValidationTarget = 'json' | 'query' | 'param' | 'form';

export function validateMiddleware<T>(target: ValidationTarget, contract: IValidator<T>) {
    return async function (ctx: AppContext, next: Next) {
        let data: any = {};

        switch (target) {
            case 'json':
                data = await ctx.req.json();
                ctx.req.addValidatedData('json', data);
                break;

            case 'query':
                data = ctx.req.query();
                ctx.req.addValidatedData('query', data);
                break;

            case 'param':
                data = ctx.req.param();
                ctx.req.addValidatedData('param', data);
                break;

            case 'form':
                data = await ctx.req.parseBody();
                ctx.req.addValidatedData('form', data);
                break;
        }

        let details = contract.validate(data);

        if (!details.success) {
            throw new ValidationError(details.errors!);
        }
        await next();
    }
}

