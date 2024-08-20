package td

func TDFilter(rows, cols uint) func(msg string) bool {
    return func(msg string) bool {
        c, err := PositionFromString(msg)
        if err != nil {
            return false
        }

        return c.Row <= rows && c.Row >= 2 &&
            c.Col <= cols && c.Col >= 2
    }
}
