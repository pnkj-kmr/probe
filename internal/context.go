package probe

import "context"

type Context struct {
	context context.Context
}

func newContext(ctx context.Context) *Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Context{
		context: ctx,
	}
}

// Context returns a context.Context, user defined on spawn or
// a context.Background as default
func (c *Context) Context() context.Context {
	return c.context
}
