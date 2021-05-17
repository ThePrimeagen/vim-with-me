import bus from "../message-bus";
import dateNow from "../now";

export default class SystemCommand {
    private stopTime: number;
    private timerId: ReturnType<typeof setInterval>;

    constructor(private onCommand: string, private offCommand: string,
                private commandLength: number) {
        this.stopTime = 0;
    }

    add() {
        const now = dateNow();
        if (now > this.stopTime) {
            bus.emit("system-command", this.onCommand);
            this.startTimer();
            this.stopTime = now;
        }

        this.stopTime += this.commandLength;
    }

    private startTimer() {
        if (this.timerId) {
            return;
        }

        this.timerId = setInterval(() => {
            const now = dateNow();
            if (now > this.stopTime) {
                clearInterval(this.timerId);
                this.timerId = null;
                bus.emit("system-command", this.offCommand);
            }
        }, 25);
    }
}

