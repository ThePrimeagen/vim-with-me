const bufferLength = 1 + 50 + 2 + 200;
const zeroBuf = Buffer.alloc(bufferLength).fill(0);

//            1             2 - 51        52 - 53     ?54 - ...?
//     +---------------+---------------+------------+------------------+
//     |     type      |   statusline  |    cost    | data ...         |
//     +---------------+---------------+------------+------------------+

export enum CommandType {
    VimCommand = 0,
    /*
    ASDF = 1,
    Xrandr = 2,
    */
    SystemCommand = 1,
    StatusUpdate = 3,
    GiveawayEnter = 4,
    VimInsert = 5,
    VimAfter = 6,
    ProgramWithMeEnter = 7,
}

const typeToString: Map<CommandType, string> = new Map([
    [CommandType.VimCommand, "VimCommand"],
    [CommandType.VimInsert, "VimInsert"],
    [CommandType.VimAfter, "VimAfter"],
    [CommandType.SystemCommand, "SystemCommand"],
    [CommandType.StatusUpdate, "StatusUpdate"],
    [CommandType.GiveawayEnter, "GiveAwayEnter"],
    [CommandType.ProgramWithMeEnter, "ProgrammWithMeEnter"],
]);

export function commandToString(type: CommandType): string {
    return typeToString.get(type);
}

const typeIdx = 0;
const statuslineIdx = 1;
const costIdx = 51;
const dataIdx = 53;

export default class Command {
    private _buffer: Buffer;
    get buffer(): Buffer {
        return this._buffer;
    }

    constructor() {
        this._buffer = Buffer.allocUnsafe(bufferLength);
    }

    reset(): Command {
        zeroBuf.copy(this._buffer);
        return this;
    }

    setType(type: CommandType): Command {
        this._buffer[typeIdx] = type;
        return this;
    }

    setStatusLine(status: string): Command {
        if (status.length > 50) {
            throw new Error("PRIME WHAT THE HELL IS GOING ON HERE? YOUR STATUS LINE THAT YOU HAVE DESIGNED IS ABOVE 50 ????????");
        }
        Buffer.from(status).copy(this._buffer, statuslineIdx);
        return this;
    }

    setCost(cost: number): Command {
        this._buffer.writeUInt16BE(cost, costIdx);
        return this;
    }

    setData(data: Buffer | null): Command {
        if (data === null) {
            return this;

        }

        if (data.length > 200) {
            throw new Error("PRIME... AGAIN????? How could you do this to me (future prime)?");
        }
        data.copy(this._buffer, dataIdx);
        return this;
    }
}

