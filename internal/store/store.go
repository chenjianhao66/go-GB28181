package store

// Factory defines the factory storage interface
type Factory interface {
	Devices() DeviceStore
}

// Store defines base storage interface
type Store[T any] interface {
	Save(entity T) error
	DeleteById(id uint) error
	Update(entity T) error
	List() ([]T, error)
	GetById(id uint) (T, error)
}

var client Factory

// Client return the store instance
func Client() Factory {
	return client
}
