package ansiparser

/*
func nextAnsiChunk(in_data []byte, idx int) (bool, int, *ansi.StyledText, error) {
	data := in_data[idx:]
	assert.Assert(data[0] == '', "the ansi chunks should always start on an escape")

	nextEsc := bytes.Index(data[1:], []byte{''}) + 1

	var styles []*ansi.StyledText = nil
	var err error = nil
	out := 0
	var complete = nextEsc != 0
	if complete {
		out = nextEsc
		styles, err = ansi.Parse(string(data[:nextEsc]))
	} else {
		styles, err = ansi.Parse(string(data))
		out = len(data)
	}

	if styles != nil && len(styles) != 0 {
		assert.Assert(len(styles) == 1, "there must only be one style at a time parsed")
		return complete, out, styles[0], err
	}

	return complete, out, nil, err
}
*/
