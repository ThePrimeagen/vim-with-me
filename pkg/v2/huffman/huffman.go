package huffman

import (
	"container/heap"
	"errors"
	"fmt"
	"strings"

	"github.com/theprimeagen/vim-with-me/pkg/v2/ascii_buffer"
	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
	byteutils "github.com/theprimeagen/vim-with-me/pkg/v2/byte_utils"
)

const HUFFMAN_ENCODE_LENGTH = 6

type huffmanNode struct {
    value int
    count int
    left *huffmanNode
    right *huffmanNode
}

func (h *huffmanNode) String() string {
    if h == nil {
        return "nil"
    }
    return fmt.Sprintf("node(%d): %d", h.count, h.value);
}

func (h *huffmanNode) debug(indent int) string {
    indentStr := strings.Repeat(" ", indent * 2)
    if h == nil {
        return fmt.Sprintf("%s-> nil\n", indentStr)
    }

    return fmt.Sprintf("%s->%s\n", indentStr, h.String()) +
        h.left.debug(indent + 1) +
        h.right.debug(indent + 1)
}

func fromValue(value int) *huffmanNode {
    return &huffmanNode{
        value: value,
        count: 1,
        left: nil,
        right: nil,
    }
}

func join(a, b *huffmanNode) *huffmanNode {
    return &huffmanNode{
        value: 0,
        count: a.count + b.count,
        left: a,
        right: b,
    }
}

func fromFreq(freq *ascii_buffer.FreqPoint) *huffmanNode {
    return &huffmanNode{
        value: freq.Val,
        count: freq.Count,
        left: nil,
        right: nil,
    }
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*huffmanNode

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].count < pq[j].count
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*huffmanNode)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

// never had this problem in my life
var HuffmanTooLarge = errors.New("huffman tree is too large")

func CalculateHuffman(freq ascii_buffer.Frequency) []byte {
    nodes := make(PriorityQueue, freq.Length(), freq.Length())
    for i, p := range freq.Points {
        nodes[i] = fromFreq(p)
    }
    heap.Init(&nodes)

    count := 1
    for len(nodes) > 1 {
        a := heap.Pop(&nodes).(*huffmanNode)
        b := heap.Pop(&nodes).(*huffmanNode)
        c := join(a, b)

        heap.Push(&nodes, c)
        count += 2
    }

    head := heap.Pop(&nodes).(*huffmanNode)

    data := make([]byte, count * HUFFMAN_ENCODE_LENGTH, count * HUFFMAN_ENCODE_LENGTH)
    encodeTree(head, data, 0)

    fmt.Println(head.debug(0))
    return data
}

func encodeTree(node *huffmanNode, data []byte, idx int) int {
    if node == nil {
        return idx
    }

    assert.Assert(idx + 2 < len(data), "idx will exceed the bounds of the huffman array during encoding")
    leftIdx := idx + HUFFMAN_ENCODE_LENGTH

    byteutils.Write16(data, idx, node.value)
    byteutils.Write16(data, idx + 2, leftIdx)

    rightIdx := encodeTree(node.left, data, leftIdx)

    byteutils.Write16(data, idx + 4, rightIdx)
    doneIdx := encodeTree(node.right, data, rightIdx)

    if leftIdx == rightIdx && leftIdx == doneIdx {
        byteutils.Write16(data, idx + 2, 0)
        byteutils.Write16(data, idx + 4, 0)
    }

    return doneIdx
}

