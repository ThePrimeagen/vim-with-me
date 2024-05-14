package iterator

type ByteIteratorResult struct {
	Done  bool
	Value int
}

type ByteIterator interface {
    Next() ByteIteratorResult
}
