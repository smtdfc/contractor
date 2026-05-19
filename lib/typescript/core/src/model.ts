import { GeneratedValidationDetails } from "./validator.js";

export interface IContractModel {
    validate(data: any): GeneratedValidationDetails;
}
