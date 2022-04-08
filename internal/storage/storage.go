package storage

type Storage[Key comparable, T any] interface {
	Put(Key, T)
	Get(Key) T
	Exist(Key) bool
}
