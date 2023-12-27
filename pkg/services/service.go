package services

import (
	"suger-clickup/pkg/clients"
	"suger-clickup/pkg/dao"
	"suger-clickup/pkg/models"
)

type Service struct {
	client clients.ClickUpClient
	dao    dao.DBDao
}

func NewService(dao dao.DBDao, client clients.ClickUpClient) *Service {
	return &Service{
		dao:    dao,
		client: client,
	}
}

func (s *Service) GetUser(email, password string) (*models.User, error) {
	user, err := s.dao.GetUser(email, password)
	return user, err
}

func (s *Service) CreateUser(email, password string) error {
	return s.dao.CreateUser(email, password)
}
