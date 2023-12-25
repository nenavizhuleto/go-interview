package core

type Repo[ID interface{}, E interface{}] interface {
	Get(ID) (E, error)
	All() ([]E, error)
	Save(E) error
	Delete(ID) error
}

type EntityRepo[E Entity] interface {
	Get(EntityID) (E, error)
	All() ([]E, error)
	Save(E) error
	Delete(EntityID) error
	Filter(func(E) bool) (E, error)
}
