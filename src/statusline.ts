import { Redemption } from "./quirk";

export default function statusLine(quirk: Redemption): string {
    return `NOT READY YET.. Stop it ${quirk.username}`;
}

