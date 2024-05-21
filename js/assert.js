/**
 * @param {any} truthy
 * @param {string} msg
 * */
export function debugAssert(truthy, msg) {
    if (!truthy) {
        debugger
        console.error(msg)
    }
}

/**
 * no idea how to make that work...
 * @param {any} truthy
 * @param {string} msg
 * @returns {asserts truthy is true}
 * */
export function assert(truthy, msg) {
    if (!truthy) {
        throw new Error(msg)
    }
}

