import { PrimeMessage } from "~/irc/prime-commands";
import { Redemption } from "~/quirk";
import { ValidationResult } from ".";

let mode: PrimeMessage = PrimeMessage.FFA;
export function primeOnlyMode(m: PrimeMessage) {
    mode = m;
}

export default function validate(data: Redemption): ValidationResult | null {
    if (mode === PrimeMessage.PrimeOnly && data.username !== "ThePrimeagen") {
        return {
            success: false,
            error: `Sorry ${data.username} its prime only mode`,
        };
    }
    return null;
}

