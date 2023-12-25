package memory

import "helpdesk/internals/core"

type Collection core.Repo[interface{}, interface{}]

type MemoryStorage struct {
	collections map[string]Collection
}
