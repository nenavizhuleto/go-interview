package core

type EntityID string
func NilEntityID() EntityID {
	return EntityID("")
}
type Children[T IChild] map[EntityID]T

type IParent interface {
	Entity
	Children() Children[IChild]
}

type IChild interface {
	Entity
	Parent() Entity
}

type Entity interface {
	EntityID() EntityID
}

type NilEntity struct {
	Entity
}

func (e NilEntity) EntityID() EntityID {
	return EntityID("")
}
