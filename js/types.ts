export type Frame = {
    cmd: number
    seq: number
    data: Uint8Array
}

export interface ByteWriter {
    write(num: number): boolean
    len(): number
    data(): Uint8Array
}

export interface SendKey {
    send(str: string): void
}
