package mycache

// getter can be a struct type or a function
type getter interface {
	Get(string) (ByteView, bool)
}

// for getters with a type of function, users should convert the function to the type GetFunc
type GetFunc func(string) (ByteView, bool)

func (f GetFunc) Get(key string) (ByteView, bool) {
	return f(key)
}
