package tcp

import (
	"math"
	"strconv"
	"strings"
)

var BASE int = 90
var ZERO int = int('$') // 36

func ToTCPInt(length int) string {
	n := length
	out := []string{}

	for ; n > 0; n = int(math.Floor(float64(n / BASE))) {
		v := n % BASE
		out = append([]string{
			strconv.Itoa(v + ZERO),
		}, out...)
	}

	return strings.Join(out, "")
}

func FromTCPInt(in_num string) int {
	num := 0
	length := len(in_num)
	for pos, val := range in_num {
		base := int(math.Pow(float64(BASE), float64(length-pos-1)))
		value := int(val) - ZERO
		num += value * base
	}

	return num
}

// --- @class ColorCompression
// --- @field table table<string>
// --- @field size number
//
// local ColorCompression = {}
// ColorCompression.__index = ColorCompression
//
// function ColorCompression:new()
//     local compression = setmetatable({
//         table = {},
//         size = 0
//     }, self)
//     return compression
// end
//
// ---@param color string
// ---@return string
// function ColorCompression:decompress(color)
//     if string.sub(color, 1, 1) == "#" then
//         table.insert(self.table, color)
//         self.size = self.size + 1
//
//         return color
//     end
//
//     -- we plus one because we are using 1 based indexing and 1 based indexing is a gift that keeps giving skill issues
//     local index = from_tcp_int(color) + 1
//     local value = self.table[index]
//
//     assert(value ~= nil, string.format("index: %d does not exist in table(%d)", index, #self.table))
//     return value
// end
//
// function ColorCompression:clear()
//     self.table = {}
//     self.size = 0
// end
//
// return {
//     to_tcp_int = to_tcp_int,
//     from_tcp_int = from_tcp_int,
//
//     ColorCompression = ColorCompression
// }
//
//

//
