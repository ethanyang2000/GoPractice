package mycache

type ByteView struct {
	b []byte
}

func (bv ByteView) Slice() []byte {
	b := make([]byte, len(bv.b))
	copy(b, bv.b)
	return b
}

func (bv ByteView) Len() int64 {
	return int64(len(bv.b))
}

func (bv ByteView) String() string {
	return string(bv.b)
}
