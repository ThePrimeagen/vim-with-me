import { EventEmitter } from "events";

type Listener = (event: string | symbol, ...args: any[]) => void;
class MessageBus extends EventEmitter {
    private toAll: Listener[];

    constructor() {
        super();
        this.toAll = [];
    }

    emit(event: string | symbol, ...args: any[]): boolean {
        this.toAll.forEach(cb => cb(event, ...args));
        return super.emit(event, ...args);
    }

    listenToAll(listener: Listener): void {
        this.toAll.push(listener);
    }

    emptyAllListener(): void {
        this.toAll = [];
    }
}

export default new MessageBus();
