package cache

import (
	"github.com/de1phin/music-transfer/internal/transfer"
)

type serviceRepository struct {
	serviceName string
	data        map[int64]interface{}
}

type cacheStorage struct {
	userState map[int64]transfer.UserState
	services  []*serviceRepository
}

func NewCacheStorage() *cacheStorage {
	cache := new(cacheStorage)
	return cache
}

func (storage *cacheStorage) AddService(serviceName string) {
	newService := serviceRepository{
		serviceName,
		make(map[int64]interface{}),
	}
	storage.services = append(storage.services, &newService)
}

func (storage *cacheStorage) PutServiceData(serviceName string, id int64, data interface{}) {

	for _, service := range storage.services {
		if service.serviceName == serviceName {
			service.data[id] = data
			break
		}
	}

}

func (storage *cacheStorage) GetServiceData(serviceName string, id int64) interface{} {

	for _, service := range storage.services {
		if service.serviceName == serviceName {
			return service.data[id]
		}
	}

	return nil
}
