import { Redemption } from "./quirk";

export default function getCost(data: Redemption): number {
    return +data.cost;
}
