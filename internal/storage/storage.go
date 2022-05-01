package storage

type Storage[Key comparable, T any] interface {
	Set(Key, T) error
	Get(Key) (T, error)
	Exist(Key) (bool, error)
}
