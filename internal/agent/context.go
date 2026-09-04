package agent

type Event struct {
	Role    string
	Content string
}

type Context struct {
	capacity int
	events   []Event
}

func NewContext(capacity int) *Context {
	return &Context{capacity: capacity}
}

func (c *Context) Add(event Event) {
	if c.capacity <= 0 {
		return
	}
	c.events = append(c.events, event)
	if len(c.events) > c.capacity {
		c.events = append([]Event(nil), c.events[len(c.events)-c.capacity:]...)
	}
}

func (c *Context) Events() []Event {
	return append([]Event(nil), c.events...)
}
