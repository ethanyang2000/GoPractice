package mycache

type ByteView struct {
	B []byte
}

func (bv ByteView) Slice() []byte {
	B := make([]byte, len(bv.B))
	copy(B, bv.B)
	return B
}

func (bv ByteView) Len() int64 {
	return int64(len(bv.B))
}

func (bv ByteView) String() string {
	return string(bv.B)
}
