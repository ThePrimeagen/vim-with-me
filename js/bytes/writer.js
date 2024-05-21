import { assert } from "../assert.js"

export class EightBitWriter {
    /** @type {Uint8Array | null} */
    #buffer = null

    /** @type {number} */
    #idx

    constructor() {
        this.#idx = 0
    }

    data() {
        return this.#buffer?.slice(0, this.#idx) ?? new Uint8Array(0)
    }

    len() {
        return this.#idx
    }

    /** @param {Uint8Array} buf */
    reset(buf) {
        this.#idx = 0
        this.#buffer = buf
    }

    /**
     * @param {number} num
     * @returns {boolean}
     **/
    write(num) {
        assert(this.#buffer !== null, "expected buffer to not equal null")

        // TODO: figure out how to make assert "type safe"
        if (this.#buffer === null) {
            return false
        }

        if (this.#buffer.byteLength === this.#idx) {
            return false
        }

        this.#buffer[this.#idx++] = num
        return true
    }
}

