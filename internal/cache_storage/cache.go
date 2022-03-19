package cache

import (
	"log"
	"sync"

	"github.com/de1phin/music-transfer/internal/transfer"
)

type serviceRepository struct {
	serviceName string
	data        map[int64]interface{}
}

type cacheStorage struct {
	mutex    sync.Mutex
	user     map[int64]transfer.User
	services []*serviceRepository
}

func NewCacheStorage() *cacheStorage {
	cache := new(cacheStorage)
	cache.user = make(map[int64]transfer.User)
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
			log.Println("Put", data)
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

func (storage *cacheStorage) HasUser(id int64) bool {
	_, ok := storage.user[id]
	return ok
}

func (storage *cacheStorage) GetUser(id int64) transfer.User {
	return storage.user[id]
}

func (storage *cacheStorage) PutUser(user transfer.User) {
	storage.mutex.Lock()
	storage.user[user.ID] = user
	log.Println("Put", user)
	storage.mutex.Unlock()
}
