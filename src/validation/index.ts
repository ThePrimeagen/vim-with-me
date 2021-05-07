import { Redemption } from "../quirk"

export type ValidationResult = {
    success: boolean;
    error?: string | null;
}

export type Validator = (data: Redemption) => ValidationResult;

const validators: Validator[] = [];

export function addValidator(validator: Validator): void {
    if (validators.includes(validator)) {
        return;
    }

    validators.push(validator);
}

export default function validate(data: Redemption): ValidationResult | null {
    let out: ValidationResult = null;

    for (let i = 0; (!out || out.success) && i < validators.length; ++i) {
        out = validators[i](data);
    }

    return out;
}

