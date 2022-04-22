package storage

type Storage[Key comparable, T any] interface {
	Put(Key, T) error
	Get(Key) (T, error)
	Exist(Key) (bool, error)
}
