import Command, { CommandType } from "../cmd";

const goldenBuffer = Buffer.from(
    //@ts-ignore
    [CommandType.ASDF].concat(
        new Array(50).fill(69)).concat([
            13, 37,
        ]).concat(new Array(200).fill(23))
);
const empty = Buffer.alloc(50 + 1 + 2 + 200).fill(0);

describe("Command", function() {

    it("should be able to set in some values.", function() {
        const cmd = new Command();
        cmd.
            setType(CommandType.ASDF).
            setStatusLine(new Array(50).fill("E").join("")).
            setCost(13 * 256 + 37).
            setData(Buffer.alloc(200).fill(23));

        expect(cmd.buffer).toEqual(goldenBuffer);
    });

    it("should reset the buffer back to 0.", function() {
        const cmd = new Command();
        cmd.
            setType(CommandType.ASDF).
            setStatusLine(new Array(50).fill("E").join("")).
            setCost(13 * 256 + 37).
            setData(Buffer.alloc(200).fill(23)).reset();
        expect(cmd.buffer).toEqual(empty);
    });

    it("should set the type", function() {
        const cmd = new Command();
        cmd.
            setType(CommandType.StatusUpdate);

        expect(cmd.buffer[0]).toEqual(CommandType.StatusUpdate);
    });
});

