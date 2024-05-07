package ascii_buffer

import (
	"container/heap"
	"errors"
	"fmt"
	"strings"

	"github.com/theprimeagen/vim-with-me/pkg/v2/assert"
)

type huffmanNode struct {
    value byte
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

func fromValue(value byte) *huffmanNode {
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

func fromFreq(freq *FreqPoint) *huffmanNode {
    return &huffmanNode{
        value: freq.idx,
        count: freq.count,
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
const HUFFMAN_ENCODE_LENGTH = 3

func CalculateHuffman(freq Frequency) ([]byte, error) {
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

    if count * HUFFMAN_ENCODE_LENGTH >= 256 {
        return nil, errors.Join(HuffmanTooLarge, fmt.Errorf("node count exceeded 255.  received %d", count))
    }

    head := heap.Pop(&nodes).(*huffmanNode)

    data := make([]byte, count * HUFFMAN_ENCODE_LENGTH, count * HUFFMAN_ENCODE_LENGTH)
    encodeTree(head, data, 0)

    return data, nil
}

func encodeTree(node *huffmanNode, data []byte, idx int) int {
    if node == nil {
        return idx
    }


    assert.Assert(idx + 2 < len(data), "idx will exceed the bounds of the huffman array during encoding")

    leftIdx := idx + HUFFMAN_ENCODE_LENGTH

    data[idx] = node.value
    data[idx + 1] = byte(leftIdx)

    rightIdx := encodeTree(node.left, data, leftIdx)

    data[idx + 2] = byte(rightIdx)
    doneIdx := encodeTree(node.right, data, rightIdx)

    if leftIdx == rightIdx && leftIdx == doneIdx {
        data[idx + 1] = 0
        data[idx + 2] = 0
    }

    return rightIdx
}
