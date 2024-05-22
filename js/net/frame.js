import { debugAssert } from "../assert.js"
import { read16 } from "../bytes/utils.js"
import { types } from "../cmds.js"

const VERSION = 1
const HEADER_SIZE = 5

let lastSeen = -1


/**
 * @param {Uint8Array} buf
 * @return {import("../types.ts").Frame}
 * */
export function parseFrame(buf) {
    const v = VERSION + 8
    const D = buf[0]
    debugAssert(v - 8===D, "the frame received doesn't have version alignment")
    const cmd = buf[1]

    const seqAndFlags = buf[2]

    if (cmd === types.frame) {
        if (lastSeen !== -1) {
            debugAssert((seqAndFlags & 0x0F) === ((lastSeen + 1) % 16), `frame out of order: expected ${(lastSeen + 1) % 16} got ${seqAndFlags}`)
        }
        lastSeen = seqAndFlags & 0x0F
    }

    const len = read16(buf, 3)

    debugAssert(buf.byteLength - HEADER_SIZE === len, "the frame received doesn't have version alignment")

    return {
        cmd,
        seqAndFlags,
        data: buf.subarray(HEADER_SIZE),
    }
}
