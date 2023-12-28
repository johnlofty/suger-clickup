package services

import (
	"errors"
	"fmt"
	"suger-clickup/pkg/clients"
	"suger-clickup/pkg/dao"
	"suger-clickup/pkg/models"
	"time"

	"github.com/gofiber/fiber/v2/log"
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

func (s *Service) CreateUser(email, password string, orgID int32) error {
	return s.dao.CreateUser(email, password, orgID)
}

func (s *Service) UpdateUserOrg(userID int32, r *models.UpdateUserRequest) error {
	org, err := s.dao.GetOrg(r.OrgID)
	if err != nil || org.OrgName == "" {
		return errors.New("invalid org_id")
	}

	user, err := s.dao.GetUserByID(userID)
	if err != nil || user.Email == "" {
		return errors.New("invalid user")
	}
	// TODO get org_id when get user
	return s.dao.UpdateUserOrg(userID, r.OrgID)
}

func (s *Service) CreateOrg(orgName string) error {
	return s.dao.CreateOrg(orgName)
}

func (s *Service) CreateTask(userID int32, r *models.CreateTaskRequest) (string, error) {
	if err := r.Validate(); err != nil {
		return "", err
	}
	user, err := s.dao.GetUserByID(userID)
	if err != nil || user.Email == "" || user.OrgId == 0 {
		return "", errors.New("invalid user")
	}
	org, err := s.dao.GetOrg(user.OrgId)
	if err != nil {
		return "", errors.New("invalid user")
	}
	task := &models.Task{
		Name:        r.Name,
		Description: r.Description,
		DueDate:     r.DueTime,
		DueDateTime: r.DueTime > 0,
	}
	task.AddCreator(fmt.Sprintf("%d", userID))
	task.AddOrgId(fmt.Sprintf("%d", org.OrgId))
	task.AddOrgName(org.OrgName)
	taskID, err := s.client.CreateTask(task)
	log.Debugf("create task taskID:%s err:%v", taskID, err)
	if err != nil {
		return "", err
	}
	ticket := &models.Ticket{
		TicketID:  taskID,
		UserID:    userID,
		OrgID:     org.OrgId,
		CreatedAt: time.Now(),
	}
	err = s.dao.CreateTicket(ticket)
	if err != nil {
		return "", fmt.Errorf("create ticket fail:%w", err)
	}
	return taskID, nil
}

func (s *Service) GetTickets(userID, page, pageSize int32) ([]models.Task, error) {
	user, err := s.dao.GetUserByID(userID)
	if err != nil || user.Email == "" || user.OrgId == 0 {
		return nil, errors.New("invalid user")
	}
	tickets, err := s.dao.GetTicketsByOrgID(user.OrgId, page, pageSize)
	if err != nil {
		return nil, err
	}
	data := make([]models.Task, 0)
	for _, ticket := range tickets {
		clickTask, err := s.client.GetTask(ticket.TicketID)
		if err != nil {
			return nil, err
		}
		task := models.Task{
			ID:          clickTask.ID,
			Name:        clickTask.Name,
			Description: clickTask.Description,
			Status:      clickTask.Status.Status,
			StartDate:   ticket.CreatedAt.Unix(),
		}
		data = append(data, task)
	}
	return data, nil
}

func (s *Service) GetTicketsCount(orgID int32) (int32, error) {
	return s.dao.GetTicketsCount(orgID)
}
