package worker

import (
	"errors"
	"sort"
	"strings"
	"sync"
)

var (
	ErrInvalidType      = errors.New("worker type is required")
	ErrDuplicateType    = errors.New("worker type is already registered")
	ErrWorkerNotFound   = errors.New("worker type is not registered")
)

// Registry owns worker registration for one Agent process. Registration is
// intentionally explicit so extensions cannot silently replace one another.
type Registry struct {
	mu      sync.RWMutex
	workers map[string]Worker
}

func NewRegistry() *Registry {
	return &Registry{workers: make(map[string]Worker)}
}

func (r *Registry) Register(w Worker) error {
	if w == nil {
		return ErrInvalidType
	}
	typeName := strings.TrimSpace(w.Type())
	if typeName == "" {
		return ErrInvalidType
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.workers[typeName]; exists {
		return ErrDuplicateType
	}
	r.workers[typeName] = w
	return nil
}

func (r *Registry) Get(typeName string) (Worker, error) {
	typeName = strings.TrimSpace(typeName)
	if typeName == "" {
		return nil, ErrInvalidType
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.workers[typeName]
	if !ok {
		return nil, ErrWorkerNotFound
	}
	return w, nil
}

// Types returns registered worker types in deterministic order.
func (r *Registry) Types() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	types := make([]string, 0, len(r.workers))
	for typeName := range r.workers {
		types = append(types, typeName)
	}
	sort.Strings(types)
	return types
}
