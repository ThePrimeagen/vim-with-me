import { assert, debugAssert } from "./assert.js"
import { types } from "./cmds.js"
import { createDecodeFrame, createOpen, expand, pushFrame } from "./decode/frame.js"
import { AsciiWindow } from "./display/window.js"
import { parseFrame } from "./net/frame.js"
import { WS } from "./ws/index.js"

// TODO: provide a url?

/**
 * @param {HTMLElement} el
 */
function run(el) {
    const wsHost = (window.location.protocol === "https:" ? "wss://" : "ws://") +
        window.location.host + ":" + window.location.port + "/ws"
    const ws = new WS(wsHost)

    /** @type {AsciiWindow | null} */
    let asciiWindow = null

    const decodeFrame = createDecodeFrame()

    ws.onMessage(async function(buf) {

        const frame = parseFrame(buf)
        switch (frame.cmd) {
        case types.open:

            if (asciiWindow !== null) {
                asciiWindow.destroy()
            }

            const open = createOpen(frame.data)
            console.log("open", open, frame.data)
            asciiWindow = new AsciiWindow(el, open.rows, open.cols / 2)

            break
        case types.frame:
            pushFrame(decodeFrame, frame)
            expand(decodeFrame)

            if (asciiWindow === null) {
                console.error("window is null?")
                return
            }

            asciiWindow.push(decodeFrame.decodeFrame)

            // render
            // console.log(data.slice(0, 5))
            break
        default:
            throw new Error("unhandled frame type " + frame.cmd)
        }

    })
}

run(document.body)
