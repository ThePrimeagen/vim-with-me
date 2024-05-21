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
 * @param {any[]} args
 * @returns {asserts truthy is true}
 * */
export function assert(truthy, msg, ...args) {
    if (!truthy) {
        console.error("assert context", args.join(" "))
        throw new Error(msg)
    }
}

