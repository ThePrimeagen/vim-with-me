/**
 * @param {any} truthy
 * @param {string} msg
 * @param {any[]} args
 * */
export function debugAssert(truthy, msg, ...args) {
    if (!truthy) {
        console.error("assert context", args.join(" "))
        console.error(msg)
        debugger
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

