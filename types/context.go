package types

type Context struct {
	store  map[string]Object
	parent *Context
}

func NewContext(parent *Context) *Context {
	return &Context{parent: parent, store: make(map[string]Object)}
}

func (ctx *Context) Set(k string, v Object) {
	ctx.store[k] = v
}

func (ctx *Context) Get(k string) (result Object, ok bool) {
	result, ok = ctx.store[k]

	if !ok && ctx.parent != nil {
		result, ok = ctx.parent.store[k]
	}

	return
}
