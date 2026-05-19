import { GeneratedValidationDetails } from "./validator.js";

export interface IContract {
    validate(data: any): GeneratedValidationDetails;
}
