import * as dotenv from "dotenv";

dotenv.config();

function getString(key: string): string {
    return process.env[key] || "";
}

function getInt(key: string): number {
    const res: string = getString(key);
    return +res || NaN;
}

export {
    getString,
    getInt
};
