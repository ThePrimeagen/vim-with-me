package encoder

import "errors"

var EncoderExceededBufferSize = errors.New("encoder exceeded buffer size of out buffer, useless encoding")

func XorRLE(frame *EncodingFrame) error {
    return nil
}

