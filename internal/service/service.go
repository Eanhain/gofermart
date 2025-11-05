package service

import (
	dto "github.com/Eanhain/gofermart/internal/api"
)

type Repository interface {
	Add(user dto.User)
	Del(user dto.User)
	List()
}

type Service struct {
	repo Repository
}

func InitialService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddUser(user dto.User) {
	s.repo.Add(user)
}
