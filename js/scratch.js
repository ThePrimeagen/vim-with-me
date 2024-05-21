import { EightBitWriter } from "./bytes/writer.js"
export const scratchBuff = new ArrayBuffer(1024 * 1024)
export const scratchArr = new Uint8Array(scratchBuff)

const writer = new EightBitWriter()
/** @returns {import("./types").ByteWriter} */
export function scratchWriter8Bit() {
    writer.reset(scratchArr)
    return writer
}

