package mycache

// getter can be a struct type or a function
type getter interface {
	Get(string) ([]byte, error)
}

// for getters with a type of function, users should convert the function to the type GetFunc
type GetFunc func(string) ([]byte, error)

func (f GetFunc) Get(key string) ([]byte, error) {
	return f(key)
}
