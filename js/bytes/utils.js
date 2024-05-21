
/**
 * @param {Uint8Array} buf
 * @param {number} offset
 */
export function read16(buf, offset) {
    return (buf[offset] << 8) + buf[offset + 1]
}

