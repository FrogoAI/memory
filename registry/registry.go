// Package registry provides a thread-safe registry; all methods are protected by a read-write mutex.
package registry

import (
	"errors"
	"sync"
	"sync/atomic"
)

const defaultBuffer = 1000

// Registry is a thread-safe collection of named groups, each holding entities keyed by ID.
type Registry[G comparable, I comparable, V any] struct {
	groups  map[G]*Group[I, V]
	indexes map[uint64]I
	aid     uint64
	mu      sync.RWMutex
}

// NewRegistry creates a new empty Registry.
func NewRegistry[G comparable, I comparable, V any]() *Registry[G, I, V] {
	return &Registry[G, I, V]{
		groups:  make(map[G]*Group[I, V]),
		indexes: make(map[uint64]I),
	}
}

// NextID atomically increments and returns the next auto-increment ID.
func (r *Registry[G, I, V]) NextID() uint64 {
	return atomic.AddUint64(&r.aid, 1)
}

// LatestID returns the current auto-increment ID without incrementing.
func (r *Registry[G, I, V]) LatestID() uint64 {
	return atomic.LoadUint64(&r.aid)
}

// SetLatestID atomically sets the auto-increment ID to the given value.
func (r *Registry[G, I, V]) SetLatestID(id uint64) {
	atomic.StoreUint64(&r.aid, id)
}

// GetGroup returns the group for the given key, creating it if it does not exist.
func (r *Registry[G, I, V]) GetGroup(key G) (group *Group[I, V]) {
	var exists bool

	r.mu.RLock()
	group, exists = r.groups[key]
	r.mu.RUnlock()

	if !exists {
		group = r.initGroup(key)
	}

	return
}

// GetGroups returns the groups for the given keys, creating any that do not exist.
func (r *Registry[G, I, V]) GetGroups(keys ...G) (groups []*Group[I, V]) {
	for _, key := range keys {
		groups = append(groups, r.GetGroup(key))
	}

	return
}

// AsyncIterator returns a channel that yields entities from the given groups concurrently.
func (r *Registry[G, I, V]) AsyncIterator(keys ...G) chan V {
	result := make(chan V, bufferSize)

	var wg sync.WaitGroup
	wg.Add(len(keys))

	for _, key := range keys {
		go func(key G) {
			defer wg.Done()

			ch := r.GetGroup(key).Iterator()

			for item := range ch {
				result <- item
			}
		}(key)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

// Iterator returns a channel that yields entities from the given groups sequentially.
func (r *Registry[G, I, V]) Iterator(keys ...G) chan V {
	result := make(chan V, bufferSize)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		for _, key := range keys {
			ch := r.GetGroup(key).Iterator()
			for item := range ch {
				result <- item
			}
		}
	}()

	go func() {
		wg.Wait()
		close(result)
	}()

	return result
}

func (r *Registry[G, I, V]) initGroup(key G) (group *Group[I, V]) {
	var exists bool

	r.mu.Lock()
	defer r.mu.Unlock()

	group, exists = r.groups[key]

	if !exists {
		group = NewGroup[I, V]()
		r.groups[key] = group
	}

	return
}

// DeleteGroup removes the group for the given key.
func (r *Registry[G, I, V]) DeleteGroup(key G) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.groups, key)
}

// AddIndex associates a uint64 ID with an entity key in the index.
func (r *Registry[G, I, V]) AddIndex(id uint64, key I) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.indexes[id] = key
}

// GetIndex returns the entity key associated with the given uint64 ID.
func (r *Registry[G, I, V]) GetIndex(id uint64) I {
	r.mu.RLock()
	defer r.mu.RUnlock()

	i := r.indexes[id]

	return i
}

// RemIndex removes the index entry for the given uint64 ID.
func (r *Registry[G, I, V]) RemIndex(id uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.indexes, id)
}

// Add inserts an entity into the specified group.
func (r *Registry[G, I, V]) Add(key G, id I, e V) error {
	group := r.GetGroup(key)
	return group.Add(id, e)
}

// Get retrieves an entity by group key and entity ID.
func (r *Registry[G, I, V]) Get(key G, id I) (e V, err error) {
	group := r.GetGroup(key)
	return group.Get(id)
}

// Remove deletes an entity from the specified group.
func (r *Registry[G, I, V]) Remove(key G, id I) error {
	group := r.GetGroup(key)
	return group.Remove(id)
}

// RemoveIDEverywhere removes the entity with the given ID from all groups.
func (r *Registry[G, I, V]) RemoveIDEverywhere(id I) error {
	var errs []error

	groups := r.GetGroups(r.GetKeys()...)

	for _, group := range groups {
		errs = append(errs, group.Remove(id))
	}

	return errors.Join(errs...)
}

// GetValues returns all entity values from the specified group.
func (r *Registry[G, I, V]) GetValues(key G) []V {
	group := r.GetGroup(key)
	return group.GetValues()
}

// TickGroup calls Tick on all entities in the specified group.
func (r *Registry[G, I, V]) TickGroup(key G) {
	group := r.GetGroup(key)
	group.Tick()
}

// TickGroups calls Tick sequentially on all entities in each specified group.
func (r *Registry[G, I, V]) TickGroups(keys ...G) {
	for _, key := range keys {
		r.TickGroup(key)
	}
}

// AsyncTick calls Tick concurrently on all entities in each specified group.
func (r *Registry[G, I, V]) AsyncTick(keys ...G) {
	var wg sync.WaitGroup
	wg.Add(len(keys))

	for _, key := range keys {
		go func(key G) {
			r.TickGroup(key)
			wg.Done()
		}(key)
	}

	wg.Wait()
}

// ClearGroup removes all entities from the specified group.
func (r *Registry[G, I, V]) ClearGroup(key G) {
	group := r.GetGroup(key)
	group.Clear()
}

// SearchInGroup searches for entities matching f within the specified group and returns results via channel.
func (r *Registry[G, I, V]) SearchInGroup(key G, f SearchFunction) chan V {
	result := make(chan V, defaultBuffer)
	group := r.GetGroup(key)

	go func(key G, result chan V, f SearchFunction) {
		group.Search(key, result, f)
		close(result)
	}(key, result, f)

	return result
}

// SearchOne returns the first entity matching f within the specified group.
func (r *Registry[G, I, V]) SearchOne(key G, f SearchFunction) V {
	group := r.GetGroup(key)
	return group.SearchOne(key, f)
}

// GetKeys returns all group keys in the registry.
func (r *Registry[G, I, V]) GetKeys() []G {
	r.mu.Lock()
	defer r.mu.Unlock()

	res := make([]G, 0, len(r.groups))
	for k := range r.groups {
		res = append(res, k)
	}

	return res
}

// Size returns the number of groups in the registry.
func (r *Registry[G, I, V]) Size() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return len(r.groups)
}
