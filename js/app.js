import { assert } from "./assert.js"
import { types } from "./cmds.js"
import { asciiPixel, createDecodeFrame, createOpen, expand, pushFrame } from "./decode/frame.js"
import { AsciiWindow } from "./display/window.js"
import { parseFrame } from "./net/frame.js"
import { WS } from "./ws/index.js"

// TODO: provide a url?

/**
 * @param {HTMLElement} el
 */
function run(el) {
    const ws = new WS("ws://localhost:8080/ws")
    /** @type {AsciiWindow | null} */
    let window = null

    const decodeFrame = createDecodeFrame()
    ws.onMessage(async function(buf) {

        const frame = parseFrame(buf)
        switch (frame.cmd) {
        case types.open:

            if (window !== null) {
                return
            }
            const open = createOpen(frame.data)
            window = new AsciiWindow(el, open.rows, open.cols / 2)
            break
        case types.frame:
            pushFrame(decodeFrame, frame)
            expand(decodeFrame)

            if (window === null) {
                console.error("window is null?")
                return
            }

            window.push(decodeFrame.decodeFrame)

            // render
            // console.log(data.slice(0, 5))
            break
        default:
            throw new Error("unhandled frame type " + frame.cmd)
        }

    })
}

run(document.body)
