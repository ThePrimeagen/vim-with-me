/** @typedef {import("./types.ts").DecodeFrame} DecodeFrame */

import { assert } from "../assert.js";
import { read16 } from "../bytes/utils.js";
import { types, encodings } from "../cmds.js";
import { scratchArr, scratchWriter8Bit } from "../scratch.js";

/**
* @param {import("../types.ts").Frame} frame
* @return DecodeFrame | null
*/
export function createFrame(frame) {
    if (frame.cmd !== types.frame) {
        return null
    }

    return { frame }
}

/**
 * @param {DecodeFrame} decode
 * @return {boolean}
*/
function isHuffmanEncoded(decode) {
    return decode.frame.data[0] === encodings.HUFFMAN
}

/**
 * @param {DecodeFrame} decode
 * @return {boolean}
*/
function isXOR_RLE(decode) {
    return decode.frame.data[0] === encodings.XOR_RLE
}

/**
 * @param {Uint8Array} decoder
 * @param {number} idx
 * @returns {number}
 **/
function left(decoder, idx) {
	assert(decoder.byteLength > idx + 5, "decoder length + idx is shorter than huffmanNode decode length")
	return read16(decoder, idx + 2)
}

/**
 * @param {Uint8Array} decoder
 * @param {number} idx
 * @returns {number}
 **/
function right(decoder, idx) {
	assert(decoder.byteLength > idx + 5, "decoder length + idx is shorter than huffmanNode decode length")
	return read16(decoder, idx + 4)
}

/**
 * @param {Uint8Array} decoder
 * @param {number} idx
 * @param {number} bit
 * @returns {number}
 **/
function jump(decoder, idx, bit) {
	if (bit === 1) {
		return right(decoder, idx)
	}
	return left(decoder, idx)
}

/**
 * @param {Uint8Array} decoder
 * @param {number} idx
 * @returns {number}
 **/
function value(decoder, idx) {
	return read16(decoder, idx)
}

/**
 * @param {Uint8Array} decoder
 * @param {number} idx
 * @returns {boolean}
 **/
function isLeaf(decoder, idx) {
	assert(decoder.byteLength > idx + 5, "decoder length + idx is shorter than huffmanNode decode length")
	return read16(decoder, idx + 2) == 0 &&
		read16(decoder, idx + 4) == 0
}

/**
 * @param {Uint8Array} decodingTree
 * @param {Uint8Array} data
 * @param {number} bitLength
 * @param {import("../types.ts").ByteWriter} writer
*/
function decodeHuffman(decodingTree, data, bitLength, writer) {
	assert(data.byteLength >= bitLength/8 + 1, "you did not provide enough data")

	let idx = 0
	let decodeIdx = 0

    outer:
	while (true) {
		for (let bitIdx = 7; bitIdx >= 0; bitIdx--) {
			const bit = (data[idx] >> bitIdx) & 0x1
			bitLength--

			decodeIdx = jump(decodingTree, decodeIdx, bit)

			if (isLeaf(decodingTree, decodeIdx)) {

				if (!writer.write(value(decodingTree, decodeIdx))) {
                    throw new Error("unable to write value into buffer")
				}

				decodeIdx = 0
			}

			if (bitLength === 0) {
				break outer
			}
		}

		idx++
	}
}

/**
 * @param {DecodeFrame} decode
*/
function expandHuffman(decode) {
    // 1 byte to encoding type (huffman)
    // 1 + 2 bytes bitLen
    // 3 + 2 bytes decodingTreeLength

    const bitLen = read16(decode.decodeFrame, 1)
    const decodingTreeLength = read16(decode.decodeFrame, 3)
    const decodingTree = decode.decodeFrame.subarray(5, 5 + decodingTreeLength)
    const writer = scratchWriter8Bit()

    decodeHuffman(decodingTree, decode.decodeFrame, bitLen, writer)
    decode.decodeFrame = writer.data()

    assert(decode.decodeFrame.byteLength > 0, "decoding failed")
}


/**
 * @param {DecodeFrame} decode
*/
function expandXOR_RLE(decode) {
    if (decode.prevDecodeFrame === null) {
        return
    }

    let idx = 0
    const data = decode.frame.data
    for (let i = 1; i < data.length; i += 2) {
        const repeat = data[i]
        const char = data[i + 1]
        for (let count = 0; count < repeat; count++, idx++) {
            scratchArr[idx] = char ^ decode.prevDecodeFrame[idx]
        }
    }

    // TODO: Copy within?
    decode.decodeFrame = scratchArr.slice(0, idx)
}

/**
 * @param {DecodeFrame} decode
*/
export function expand(decode) {
    if (isXOR_RLE(decode)) {
        expandXOR_RLE(decode)
    } else {
        expandHuffman(decode)
    }
}

/**
 * @param {DecodeFrame} decode
 * @return {Uint8Array}
*/
export function asciiPixel(decode) {
    const frame = decode.decodeFrame
    const out = new Uint8Array(frame.byteLength)
    for (let j = 0, i = 0; i < frame.byteLength; ++i, j += 2) {
        out[j] = frame[i]
        out[j + 1] = frame[i]
    }
    return out
}


