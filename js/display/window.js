
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
 * @returns {string}
 */
function toRgb(color) {
    const r = (color & redMask) >> 5
    const g = (color & greenMask) >> 3
    const b = (color & blueMask)

    return `#${toHex(r * 255 / 7)}${toHex(g * 255 / 3)}${toHex(b * 255 / 7)}`
}

export class AsciiWindow {
    /** @type {HTMLElement} */
    #el

    /** @type {HTMLDivElement} */
    #container

    /** @type {HTMLDivElement[][]} */
    #divs

    /** @type {number[][]} */
    #colors

    /** @type {number} */
    #rows
    /** @type {number} */
    #cols

    /**
     * @param {HTMLElement} el
     * @param {number} rows
     * @param {number} cols
     * */
    constructor(el, rows, cols) {
        this.#colors = []
        this.#divs = []
        this.#el = el
        this.#container = document.createElement("div")
        this.#rows = rows
        this.#cols = cols

        this.#init()
    }

    /**
     * @param {Uint8Array} frame
     */
    push(frame) {
        let count = 0
        let black = 0

        for (let row = 0; row < this.#rows; ++row) {
            for (let col = 0; col < this.#cols; ++col) {
                const idx = row * this.#cols + col
                const color = this.#colors[row][col]
                const inColor = frame[idx]

                if (color === inColor) {
                    continue
                }

                count++
                this.#divs[row][col].style.backgroundColor = toRgb(inColor)
                this.#colors[row][col] = inColor
            }
        }
    }

    #init() {
        /** @type {HTMLDivElement[][]} */
        const divs = this.#divs
        const cont = this.#container
        const colors = this.#colors

        cont.style.display = "grid"
        cont.style.gridTemplateColumns = `repeat(${this.#cols}, 1fr)`
        cont.style.gridTemplateRows = `repeat(${this.#rows}, 1fr)`
        cont.style.height = "100vh"

        for (let row = 0; row < this.#rows; ++row) {
            divs[row] = []
            colors[row] = []

            for (let col = 0; col < this.#cols; ++col) {
                const div = document.createElement("div")
                divs[row].push(div)

                div.style.backgroundColor = "#000"
                colors[row].push(0)

                cont.appendChild(div)
            }
        }

        this.#el.appendChild(cont)
    }

    destroy() {
        this.#el.innerHTML = ""
    }
}

