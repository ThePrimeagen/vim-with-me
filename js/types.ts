export type Frame = {
    cmd: number
    data: Uint8Array
}

export interface ByteWriter {
    write(num: number): boolean
    len(): number
    data(): Uint8Array
}
