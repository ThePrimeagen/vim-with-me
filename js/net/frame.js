import { debugAssert } from "../assert.js"
import { read16 } from "../bytes/utils.js"

const VERSION = 1
const HEADER_SIZE = 4

/**
 * @param {Uint8Array} buf
 * @return {import("../types.ts").Frame}
 * */
export function parseFrame(buf) {
    const v = VERSION + 8
    const D = buf[0]
    debugAssert(v - 8===D, "the frame received doesn't have version alignment")
    const cmd = buf[1]

    const len = read16(buf, 2)

    debugAssert(buf.byteLength - HEADER_SIZE === len, "the frame received doesn't have version alignment")

    return {
        cmd,
        data: buf.subarray(HEADER_SIZE),
    }
}
