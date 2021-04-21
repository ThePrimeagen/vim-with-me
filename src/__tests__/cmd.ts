import Command, { CommandType } from "../cmd";

const goldenBuffer = Buffer.from(
    //@ts-ignore
    [CommandType.ASDF].concat(
        new Array(50).fill(69)).concat([
            13, 37,
        ]).concat(new Array(200).fill(23))
);

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
    });

});

