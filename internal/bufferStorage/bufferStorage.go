package bufferstorage

import (
	dto "github.com/Eanhain/gofermart/internal/api"
	cache "github.com/patrickmn/go-cache"
)

type Storage interface {
	Add(dto.User) error
	Del(dto.User) error
	MultipleAdd(dto.UserArray) error
	MultipleDel(dto.UserArray) error
	List() ([]dto.UserArray, error)
}

type BufferStorage struct {
	storage Storage
	cache   *cache.Cache
}

func InitialBufferStorage(storage Storage) *BufferStorage {
	c := cache.New(0, 0)
	return &BufferStorage{
		storage: storage,
		cache:   c}
}

func (buffer *BufferStorage) Add(info dto.User) error {
	err := buffer.storage.Add(info)
	return err
}

func (buffer *BufferStorage) Del(info dto.User) error {
	err := buffer.storage.Del(info)
	return err
}

func (buffer *BufferStorage) List(info dto.User) ([]dto.UserArray, error) {
	arr, err := buffer.storage.List()
	return arr, err
}
