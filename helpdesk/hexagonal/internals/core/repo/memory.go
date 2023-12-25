package repo

import (
	"errors"
	"helpdesk/internals/core"
)

var (
	ErrEntityNotFound = errors.New("entity not found")
)

type MemoryRepo[E core.Entity] struct {
	name string
	c    map[core.EntityID]E
}

func NewMemoryRepo[E core.Entity](name string) MemoryRepo[E] {
	return MemoryRepo[E]{
		name: name,
		c:    make(map[core.EntityID]E),
	}
}

func (r MemoryRepo[E]) Get(id core.EntityID) (E, error) {
	if e, ok := r.c[id]; ok {
		return e, nil
	} else {
		var n E
		return n, ErrEntityNotFound
	}
}

func (r MemoryRepo[E]) All() ([]E, error) {
	entities := make([]E, 0)
	for _, e := range r.c {
		entities = append(entities, e)
	}

	return entities, nil
}

func (r *MemoryRepo[E]) Save(e E) error {
	r.c[e.EntityID()] = e
	return nil
}

func (r *MemoryRepo[E]) Delete(id core.EntityID) error {
	if _, ok := r.c[id]; ok {
		delete(r.c, id)
		return nil
	} else {
		return ErrEntityNotFound
	}
}

func (r *MemoryRepo[E]) Filter(f func(E) bool) (E, error) {
	for _, c := range r.c {
		if f(c) {
			return c, nil
		}
	}
	var e E
	return e, ErrEntityNotFound
}
