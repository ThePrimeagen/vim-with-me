package td

type BoxPos struct {
    maxRows int
    position int
}

func NewBoxPos(maxRows int) BoxPos {
    return BoxPos{
        maxRows: maxRows,
        position: 0,
    }
}

// As of right now, out of bounds guesses will place a tower within your area
// randomly
func (r *BoxPos) NextPos() Position {
    col := 6
    if r.position & 0x1 == 0 {
        col = 12
    }

    row := ((r.position % 4) / 2) * 5
    r.position++

    return Position{
        Row: uint(row),
        Col: uint(col),
    }
}

