import { Redemption } from "../quirk";
import { Validator, ValidationResult } from "../validation";
import bus from "../message-bus";
import getType from "../get-type";
import { CommandType } from "../cmd";

export default class ProgramWithMe {

    public banned: string[] = [];
    public programmers: string[] = [];
    public enabled = false;
    public currentFamousPerson = "";
    public countdownId: ReturnType<typeof setTimeout>;
    public validateFunction: Validator;
    private famousIdx: number;

    constructor(public timeoutTime = 25000) {
        this.validateFunction = this.validateProgramWithMe.bind(this);

        bus.on("quirk-message", (data: Redemption) => {
            const type = getType(data);
            if (type === CommandType.ProgramWithMeEnter) {
                this.addToProgramWithMe(data);
            }
        });

    }

    private validateProgramWithMe(data: Redemption): ValidationResult {
        if (!this.enabled) {
            return { success: true };
        }

        const type = getType(data);

        if (type !== CommandType.VimAfter &&
            type !== CommandType.VimInsert &&
                type !== CommandType.VimCommand) {
            return { success: true };
        }

        if (this.currentFamousPerson === "") {
            this.setNextFamousPerson();
        }

        if (data.username === this.currentFamousPerson) {
            this.setNextFamousPerson();
            return { success: true };
        }

        return {
            success: false,
            error: `You are not the famous person ${data.username}`,
        };
    }

    enableProgramWithMe(): void {
        this.enabled = true;
        this.famousIdx = 0;
        this.setFamousPersonOrder();
        this.setNextFamousPerson();
    }

    disableProgramWithMe(): void {
        this.enabled = false;
        this.programmers = [];
        this.currentFamousPerson = "";
        this.cleanPp();
    }

    addToProgramWithMe(data: Redemption): void {
        if (!this.programmers.includes(data.username)) {
            this.programmers.push(data.username);
        }
    }

    addBannedList(bans: string[]): void {
        this.banned = this.banned.concat(bans);
    }

    private setFamousPersonOrder() {
        this.programmers.sort(() => {
            return Math.floor(Math.random() * 2) === 0 ? -1 : 1
        });
    }

    private setNextFamousPerson() {
        this.cleanPp();

        const idx = this.famousIdx++ % this.programmers.length;
        const nextIdx = this.famousIdx % this.programmers.length;

        this.currentFamousPerson = this.programmers[idx];
        this.countdownId = setTimeout(() => {
            this.setNextFamousPerson();
        }, this.timeoutTime);

        bus.emit("irc-message", `@${this.currentFamousPerson} please do a vim insert, after, or command.`);
        bus.emit("irc-message", `@${this.programmers[nextIdx]} You are next! You probably want VimAfter`);
    }

    private cleanPp() {
        if (this.countdownId) {
            clearTimeout(this.countdownId);
            this.countdownId = null;
        }
    }
}


