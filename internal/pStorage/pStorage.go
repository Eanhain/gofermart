package pStorage

import (
	dto "github.com/Eanhain/gofermart/internal/api"
)

type PersistStorage struct {
}

func InitialPersistStorage() *PersistStorage {
	return &PersistStorage{}
}

func (ps *PersistStorage) Add(user dto.User) error {
	return nil
}

func (ps *PersistStorage) Del(user dto.User) error {
	return nil
}

func (ps *PersistStorage) MultipleAdd(user dto.UserArray) error {
	return nil
}

func (ps *PersistStorage) MultipleDel(user dto.UserArray) error {
	return nil
}

func (ps *PersistStorage) List() ([]dto.UserArray, error) {
	return []dto.UserArray{}, nil
}
