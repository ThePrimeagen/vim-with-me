import { types } from "./cmds.js"
import { createDecodeFrame, createOpen, expand, pushFrame } from "./decode/frame.js"
import { AsciiWindow } from "./display/window.js"
import { Cache } from "./net/cache.js"
import { parseFrame } from "./net/frame.js"
import { WS } from "./ws/index.js"

// TODO: provide a url?

/**
 * @param {HTMLElement} el
 */
function run(el) {
    // Note: host contains port
    const wsHost = (window.location.protocol === "https:" ? "wss://" : "ws://") +
        window.location.host + "/ws"
    const ws = new WS(wsHost)

    /** @type {AsciiWindow | null} */
    let asciiWindow = null

    const decodeFrame = createDecodeFrame()
    const cache = new Cache()

    function render() {
        cache.seek()
        let f = cache.pop()
        if (f) {
            pushFrame(decodeFrame, f)
            expand(decodeFrame)
            if (asciiWindow === null) {
                console.error("window is null?")
                return
            }

            asciiWindow.push(decodeFrame.decodeFrame)
        }

        requestAnimationFrame(render)
    }
    requestAnimationFrame(render)

    ws.onMessage(async function(buf) {

        const frame = parseFrame(buf)
        switch (frame.cmd) {
        case types.open:

            if (asciiWindow !== null) {
                asciiWindow.destroy()
            }

            cache.reset()
            const open = createOpen(frame.data)
            console.log("open", open, frame.data)
            asciiWindow = new AsciiWindow(el, open.rows, open.cols / 2)

            break
        case types.frame:
            cache.push(frame)
            break
        default:
            throw new Error("unhandled frame type " + frame.cmd)
        }

    })
}

run(document.body)
