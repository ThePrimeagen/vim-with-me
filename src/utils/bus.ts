import bus from "../message-bus";

export default function accumulate() {
    const out: any[][] = [];
    bus.listenToAll(function(...args: any[]) {
        out.push(args);
    });

    return function flush(stop?: boolean) {
        if (stop) {
            bus.emptyAllListener();
        }
        return out;
    };
}

