package session

import "errors"

var ErrInvalidHistoryCapacity = errors.New("session history capacity must be greater than zero")

type Action struct {
	Tool    string
	Summary string
}

type History struct {
	capacity int
	actions  []Action
}

func NewHistory(capacity int) (*History, error) {
	if capacity <= 0 {
		return nil, ErrInvalidHistoryCapacity
	}
	return &History{
		capacity: capacity,
		actions:  make([]Action, 0, capacity),
	}, nil
}

func (h *History) Add(action Action) {
	if len(h.actions) == h.capacity {
		copy(h.actions, h.actions[1:])
		h.actions[len(h.actions)-1] = action
		return
	}
	h.actions = append(h.actions, action)
}

func (h *History) Actions() []Action {
	return append([]Action(nil), h.actions...)
}
