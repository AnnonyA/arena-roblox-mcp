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

func (c *Context) Clear() {
	c.events = nil
}

func (c *Context) CompactToolOutputs(maxRunes int) {
	if maxRunes < 0 || len(c.events) < 2 {
		return
	}
	for i := 0; i < len(c.events)-1; i++ {
		if c.events[i].Role != "tool" {
			continue
		}
		runes := []rune(c.events[i].Content)
		if len(runes) <= maxRunes {
			continue
		}
		c.events[i].Content = string(runes[:maxRunes]) + "… [compacted]"
	}
}

func (c *Context) Events() []Event {
	return append([]Event(nil), c.events...)
}
