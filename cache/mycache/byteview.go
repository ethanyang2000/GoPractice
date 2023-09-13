package mycache

type ByteView struct {
	b []byte
}

func (bv ByteView) Slice() []byte {
	B := make([]byte, len(bv.b))
	copy(B, bv.b)
	return B
}

func (bv ByteView) Len() int64 {
	return int64(len(bv.b))
}

func (bv ByteView) String() string {
	return string(bv.b)
}
