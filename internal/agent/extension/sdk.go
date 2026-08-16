package extension

import "context"

// Device identifies the target device for an extension operation.
type Device struct {
	ID       string
	Protocol string
	Address  string
}

// Request is the opaque input passed from a worker to an extension.
type Request struct {
	Device  Device
	Action  string
	Payload []byte
}

// Response is the normalized result returned by an extension.
type Response struct {
	Payload []byte
}

// Extension is the stable contract implemented by GenieACS integrations.
type Extension interface {
	Name() string
	Execute(context.Context, Request) (Response, error)
}

// Registry resolves extensions by name.
type Registry struct {
	items map[string]Extension
}

func NewRegistry(items ...Extension) *Registry {
	r := &Registry{items: make(map[string]Extension, len(items))}
	for _, item := range items {
		if item != nil {
			r.items[item.Name()] = item
		}
	}
	return r
}

func (r *Registry) Get(name string) (Extension, bool) {
	item, ok := r.items[name]
	return item, ok
}
