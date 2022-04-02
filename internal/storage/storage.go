package storage

type Storage[T any] interface {
	Put(int64, T)
	Get(int64) T
	Exist(int64) bool
}
