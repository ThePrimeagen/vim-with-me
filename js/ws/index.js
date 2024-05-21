const CONNECTING = 0
const CONNECTED = 1
const ERROR = 2
const CLOSE = 3

/** @typedef {(buf: Blob) => void} OnMessage */

class WS {
    /** @type {WebSocket} */
    #ws

    /** @type {string} */
    #url

    /** @type {number} */
    #state

    /** @type {null | OnMessage} */
    #onmessage = null

    /**
    * @param {string} url
    **/
    constructor(url) {
        this.#url = url
        this.#connect()
        this.#onmessage = null
    }

    /** @param {OnMessage} msg */
    onMessage(msg) {
        this.#onmessage = msg
    }

    push() {
        throw new Error("cannot send messages up to the server as of right now, fuck off")
    }

    #connect() {
        const ws = this.ws = new WebSocket(this.#url)
        this.#state = CONNECTING

        ws.onopen = () => {
            this.#state = CONNECTED
        }

        ws.onclose = () => {
            this.#state = CLOSE
            // some reporting or something??
            // some backoff?
            this.#connect()
        }

        ws.onerror = () => {
            this.#state = ERROR
            // some reporting or something??
            // some backoff?
            this.#connect()
        }

        ws.onmessage = (msg) => {
            /** @type {ArrayBuffer} */
            const arrBuff = msg.data

            if (this.#onmessage) {
                this.#onmessage(arrBuff);
            }
        }
    }
}


export {
    WS
}
