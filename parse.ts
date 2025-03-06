import fs from "fs";
import parse from "csv-simple-parser";

const file = Bun.file(process.argv[2]);
type Rec = {
    won: string
    round: number
    "one-team": string,
    "two-team": string,
    seed: number,
    oneTotalTowersBuild: number,
    oneTotalProjectiles: number,
    oneTotalTowerUpgrades: number,
    oneTotalCreepDamage: number,
    oneTotalTowerDamage: number,
    oneTotalDamageFromCreeps: number,
    twoTotalTowersBuild: number,
    twoTotalProjectiles: number,
    twoTotalTowerUpgrades: number,
    twoTotalCreepDamage: number,
    twoTotalTowerDamage: number,
    twoTotalDamageFromCreeps: number
};
let csv = parse(await file.text(), { header: true }) as Rec[];

csv = csv.filter(r => r.won !== "").map(r => {
    r["one-team"] = r["one-team"].split(":")[2]
    r["two-team"] = r["two-team"].split(":")[2]
    if (r.won === "49") {
        r.won = r["one-team"]
    } else {
        r.won = r["two-team"]
    }
    return r
});

const GPT4_5 = "gpt-4.5-preview"
const GPTO3 = "o3-mini"
const GPT4 = "gpt-4"

const gpt4_5 = csv.filter(r => r["one-team"] === GPT4_5 || r["two-team"] === GPT4_5)
const gpto3 = csv.filter(r => r["one-team"] === GPTO3 || r["two-team"] === GPTO3)
const gpt4 = csv.filter(r => r["one-team"] === GPT4 || r["two-team"] === GPT4)
const header = "won,round,one-team,two-team,seed,oneTotalTowersBuild,oneTotalProjectiles,oneTotalTowerUpgrades,oneTotalCreepDamage,oneTotalTowerDamage,oneTotalDamageFromCreeps,twoTotalTowersBuild,twoTotalProjectiles,twoTotalTowerUpgrades,twoTotalCreepDamage,twoTotalTowerDamage,twoTotalDamageFromCreeps"

fs.writeFileSync("gpt4_5.csv", header + "\n" + gpt4_5.map(r => Object.values(r).join(",")).join("\n"))
fs.writeFileSync("gpt4.csv", header + "\n" + gpt4.map(r => Object.values(r).join(",")).join("\n"))
fs.writeFileSync("gpto3.csv", header + "\n" + gpto3.map(r => Object.values(r).join(",")).join("\n"))

console.log(csv)

