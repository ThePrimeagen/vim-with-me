import { debugAssert } from "../assert.js"
import { isHuffmanEncoded } from "../decode/frame.js"

export class Cache {
    /** @type {(import("../types").Frame | undefined)[]} */
    cache

    /** @type {number} */
    cacheIdx

    constructor() {
        this.cache = new Array(256).fill(undefined)
        this.cacheIdx = -1
    }

    reset() {
        this.cache = new Array(256).fill(undefined)
        this.cacheIdx = -1
    }

    /** @param {import("../types").Frame} f */
    push(f) {
        debugAssert(this.cache[f.seq] === undefined, "cache collision of frame", "frame", f)
        this.cache[f.seq] = f

        if (this.cacheIdx === -1) {
            this.cacheIdx = f.seq
        }
    }

    /** @returns {import("../types").Frame | undefined} */
    pop() {
        if (this.cache[this.cacheIdx] === undefined) {
            return
        }

        const item = this.cache[this.cacheIdx]
        this.cache[this.cacheIdx] = undefined
        this.cacheIdx = (this.cacheIdx + 1) % 256

        return item
    }

    seek() {
        let idx = (this.cacheIdx + 1) % 256
        let item = this.cache[idx]
        while (item) {

            if (isHuffmanEncoded(item.data)) {
                for (let i = this.cacheIdx; i !== idx; i = (i + 1) % 256) {
                    this.cache[i] = undefined
                }
                this.cacheIdx = idx
                return
            }

            idx = (idx + 1) % 256
            item = this.cache[idx]
        }
    }
}

