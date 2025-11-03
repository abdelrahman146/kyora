package ctxkey

type ContextKey struct {
	Name string
}

func (k ContextKey) String() string {
	return "ctxkey:" + k.Name
}

func New(name string) ContextKey {
	return ContextKey{
		Name: name,
	}
}
