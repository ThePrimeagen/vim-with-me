package rgb

const NONE = 0
const RLE_ONLY = 1

type BufferEncoding int

func (b BufferEncoding) NoEncoding() BufferEncoding {
	return NONE
}

func (b BufferEncoding) RleOnly() BufferEncoding {
	return RLE_ONLY
}
