import assert from "node:assert"
import { test } from "node:test"
import { createDecodeFrame, decodeHuffman, expand, pushFrame } from "./frame.js"
import { scratchWriter8Bit } from "../scratch.js"
import { encodings, types } from "../cmds.js"

const decodeTree = [0, 0, 0, 6, 0, 12, 0, 65, 0, 0, 0, 0, 0, 0, 0, 18, 0, 24, 0, 66, 0, 0, 0, 0, 0, 0, 0, 30, 0, 36, 0, 68, 0, 0, 0, 0, 0, 67, 0, 0, 0, 0]
const A = "A".charCodeAt(0)
const C = "C".charCodeAt(0)
const expected = new Uint8Array([
    A, A, A, A, A, A,
    C,
])

test("frame decoding", function() {
    const decode = createDecodeFrame()
    const bitLen = 9
    const frame = {
        cmd: types.frame,
        data: new Uint8Array([
            // encoding
            encodings.HUFFMAN,

            // bitlen
            0,
            bitLen,

            // decodeTree Length
            0,
            decodeTree.length,

            // decodeTree
            ...decodeTree,

            // decodeable data
            3, 128
        ])
    }

    pushFrame(decode, frame)
    expand(decode)
    assert.deepEqual(expected, decode.decodeFrame, "expected vs writer.data")
})

test("huffman decoding", function() {
    const bitLen = 9
    const decodeTree = new Uint8Array([
        0, 0, 0, 6, 0, 12, 0, 65, 0, 0, 0, 0, 0, 0, 0, 18, 0, 24, 0, 66, 0, 0, 0, 0, 0, 0, 0, 30, 0, 36, 0, 68, 0, 0, 0, 0, 0, 67, 0, 0, 0, 0
    ])
    const data = new Uint8Array([3, 128])
    const writer = scratchWriter8Bit()

    decodeHuffman(decodeTree, data, bitLen, writer)

    assert.deepEqual(expected, writer.data(), "expected vs writer.data")
})
