import { assert } from "../assert.js"

/**
 * @param {number} val
 * @returns {string}
 */
function toHex(val) {
    val = Math.round(val)
    let out = val.toString(16)
    if (out.length === 1) {
        return "0" + out
    }

    return out
}

const redMask = 0b111_00_000
const greenMask = 0b000_110_00
const blueMask = 0b000_00_111

/**
 * @param {number} color
 * @returns {Array<number>}
 */
function toRgb(color) {
    const r = (color & redMask) >> 5
    const g = (color & greenMask) >> 3
    const b = (color & blueMask)

    return [Math.round(r * 255 / 7), Math.round(g * 255 / 3), Math.round(b * 255 / 7)]
}

export class AsciiWindow {
    /** @type {HTMLElement} */
    #el

    /** @type {HTMLCanvasElement} */
    #canvas

    /** @type {number[][]} */
    #colors

    /** @type {number} */
    #rows
    /** @type {number} */
    #cols

    /** @type {?CanvasRenderingContext2D} */
    #context

    /**
     * @param {HTMLElement} el
     * @param {number} rows
     * @param {number} cols
     * */
    constructor(el, rows, cols) {
        this.#colors = []
        this.#el = el
        this.#canvas = document.createElement("canvas")
        this.#rows = rows
        this.#cols = cols
        this.#context = null

        this.#init()
    }

    /**
     * @param {Uint8Array} frame
     */
    push(frame) {
        if(!this.#context){
            throw new Error("No context")
        }
        const imageData = this.#context.getImageData(0, 0, this.#canvas.width, this.#canvas.height);

        const data = imageData.data;

        for (let row = 0; row < this.#rows; ++row) {
            for (let col = 0; col < this.#cols; ++col) {

                const idx = row * this.#cols + col
                const color = this.#colors[row][col]
                const inColor = frame[idx]

                if (color === inColor) {
                    continue
                }

                const pixel = toRgb(inColor)

                data[4 * idx] = pixel[0]
                data[4 * idx + 1] = pixel[1]
                data[4 * idx + 2] = pixel[2]

                this.#colors[row][col] = inColor
            }
        }
        this.#context.putImageData(imageData, 0, 0);
    }

    /** 
    * @param {MouseEvent} event
    * @param {number} scale
    */
    #scaleCanvas(event, scale) {
        this.#canvas.style.width  = (this.#cols * scale) + "px"
        this.#canvas.style.height = (this.#rows * scale) + "px"

        for ( let btn of document.getElementsByTagName('button') ){
            btn.style.fontWeight = 'normal';
        }

        event.target.style.fontWeight = 'bold'

    }

    #init() {
        const colors = this.#colors

        this.#canvas.width  = this.#cols
        this.#canvas.height = this.#rows
        this.#canvas.style.width  = this.#cols + "px"
        this.#canvas.style.height = this.#rows + "px"

        for (let row = 0; row < this.#rows; ++row) {
            colors[row] = []

            for (let col = 0; col < this.#cols; ++col) {
                colors[row].push(0)
            }
        }

        
        const btnDiv = document.createElement('div')
        this.#el.appendChild(btnDiv)

        ;[1, 2, 4, 8].forEach((v) => {
            const btn = document.createElement('button')
            btn.textContent = "x" + v;
            btn.addEventListener("click", (b)=>this.#scaleCanvas(b, v));
            btn.style.fontSize = "110%"
            btnDiv.appendChild(btn)
        })
                

        this.#el.appendChild(this.#canvas)

        this.#context = this.#canvas.getContext('2d')
        if (this.#context !== null) {
            this.#context.fillStyle = "black";
            this.#context.fillRect(0, 0, this.#canvas.width, this.#canvas.height)
        } else {
            throw new Error("Can't load canvas context!");
        }
    }

    destroy() {
        this.#el.innerHTML = ""
    }

}

