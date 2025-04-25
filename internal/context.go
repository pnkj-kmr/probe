package probe

import (
	"context"
)

type Context struct {
	context  context.Context
	event    *EventEngine
	api      *APIEngine
	db       *DBEngine
	poll     *PollEngine
	schedule *ScheduleEngine
	export   *ExportEngine
}

func newContext(c context.Context) *Context {
	if c == nil {
		c = context.Background()
	}
	return &Context{
		context: c,
	}
}

// Context returns a context.Context, user defined on spawn or
// a context.Background as default
func (c *Context) Context() context.Context {
	return c.context
}

func (c *Context) WithEvent(e *EventEngine) *Context {
	c.event = e
	return c
}

func (c *Context) Event() *EventEngine {
	return c.event
}

func (c *Context) WithPoll(e *PollEngine) *Context {
	c.poll = e
	return c
}

func (c *Context) Poll() *PollEngine {
	return c.poll
}

func (c *Context) WithDB(e *DBEngine) *Context {
	c.db = e
	return c
}

func (c *Context) DB() *DBEngine {
	return c.db
}

func (c *Context) WithSchdule(e *ScheduleEngine) *Context {
	c.schedule = e
	return c
}

func (c *Context) Schedule() *ScheduleEngine {
	return c.schedule
}

func (c *Context) WithExport(e *ExportEngine) *Context {
	c.export = e
	return c
}

func (c *Context) Export() *ExportEngine {
	return c.export
}

func (c *Context) WithAPI(e *APIEngine) *Context {
	c.api = e
	return c
}

func (c *Context) API() *APIEngine {
	return c.api
}
