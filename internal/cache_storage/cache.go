package cache

import (
	"sync"

	"github.com/de1phin/music-transfer/internal/transfer"
)

type serviceRepository struct {
	serviceName string
	data        map[int64]interface{}
}

type cacheStorage struct {
	mutex     sync.Mutex
	userState map[int64]transfer.UserState
	services  []*serviceRepository
}

func NewCacheStorage() *cacheStorage {
	cache := new(cacheStorage)
	cache.userState = make(map[int64]transfer.UserState)
	return cache
}

func (storage *cacheStorage) AddService(serviceName string) {
	newService := serviceRepository{
		serviceName,
		make(map[int64]interface{}),
	}
	storage.services = append(storage.services, &newService)
}

func (storage *cacheStorage) PutServiceData(id int64, serviceName string, data interface{}) {

	for _, service := range storage.services {
		if service.serviceName == serviceName {
			storage.mutex.Lock()
			service.data[id] = data
			storage.mutex.Unlock()
			break
		}
	}

}

func (storage *cacheStorage) GetServiceData(id int64, serviceName string) interface{} {

	for _, service := range storage.services {
		if service.serviceName == serviceName {
			return service.data[id]
		}
	}

	return nil
}

func (storage *cacheStorage) GetUserState(id int64) transfer.UserState {
	return storage.userState[id]
}

func (storage *cacheStorage) PutUserState(id int64, userState transfer.UserState) {
	storage.mutex.Lock()
	storage.userState[id] = userState
	storage.mutex.Unlock()
}
