import { debugAssert } from "../assert.js"
import { read16 } from "../bytes/utils.js"
import { types } from "../cmds.js"

const VERSION = 1
const HEADER_SIZE = 6

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
    const seq = buf[2]
    //const flags = buf[3]

    if (cmd === types.open) {
        lastSeen = -1
    } else if (cmd === types.frame) {
        if (lastSeen !== -1) {
            debugAssert(seq === ((lastSeen + 1) % 256), `frame out of order: expected ${(lastSeen + 1) % 16} got ${seq}`)
        }
        lastSeen = seq
    }

    const len = read16(buf, 4)

    debugAssert(buf.byteLength - HEADER_SIZE === len, "the frame received doesn't have version alignment")

    return {
        cmd,
        seq: seq,
        data: buf.subarray(HEADER_SIZE),
    }
}
