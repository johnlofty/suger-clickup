package services

import (
	"errors"
	"fmt"
	"strconv"
	"suger-clickup/pkg/clients"
	"suger-clickup/pkg/dao"
	"suger-clickup/pkg/models"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

const TicketCommentPageSize = 25

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

func (s *Service) ModOrgNotification(orgID int32,
	req *models.OrgNotiRequest) error {
	return s.dao.ModOrgNotification(orgID, req)
}

func (s *Service) GetOrgNotification(orgID int32) (*models.OrgNotiRequest, error) {
	return s.dao.GetOrgNotification(orgID)
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
		Name:          r.Name,
		Description:   r.Description,
		DueDate:       r.DueTime,
		DueDateTime:   true,
		StartDate:     r.StartTime,
		StartDateTime: true,
		Priority:      r.Priority,
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
		Watcher:   r.Watcher,
	}
	err = s.dao.CreateTicket(ticket)
	if err != nil {
		return "", fmt.Errorf("create ticket fail:%w", err)
	}
	return taskID, nil
}

func (s *Service) FormatTask(clickTask models.ClickupTask, ticket models.Ticket) models.Task {
	task := models.Task{
		ID:          clickTask.ID,
		Name:        clickTask.Name,
		Description: clickTask.Description,
		Status:      clickTask.Status.Status,
		StartDate:   ticket.CreatedAt.Unix(),
		Watcher:     ticket.Watcher,
	}
	dueDate, err := strconv.ParseInt(clickTask.DueDate, 10, 64)
	if err == nil && dueDate != 0 {
		task.DueDate = dueDate / 1000
	}
	priority, err := strconv.ParseInt(clickTask.Priority.ID, 10, 64)
	if err == nil {
		task.Priority = int32(priority)
	}
	return task
}

func (s *Service) GetTickets(userID, page, pageSize int32) ([]models.Task, error) {
	user, err := s.dao.GetUserByID(userID)
	if err != nil || user.Email == "" || user.OrgId == 0 {
		return nil, errors.New("invalid user")
	}
	tickets, err := s.dao.GetTicketsByOrgID(user.OrgId, page, pageSize)
	log.Debugf("tickets:%+v err:%v", tickets, err)
	if err != nil {
		return nil, err
	}
	data := make([]models.Task, 0)
	for _, ticket := range tickets {
		clickTask, err := s.client.GetTask(ticket.TicketID)
		log.Debugf("clickTask:%v err:%v", clickTask, err)
		if err != nil {
			if err == clients.ErrTaskNotFound {
				err = s.dao.DeleteTicket(ticket.TicketID)
				log.Warnf("delete ticket:%s err:%v", ticket.TicketID, err)
			}
			continue
		}
		task := s.FormatTask(clickTask, ticket)
		data = append(data, task)
	}
	return data, nil
}

func (s *Service) GetTicket(ticketID string) (models.Task, error) {
	var task models.Task
	ticket, err := s.dao.GetTicket(ticketID)
	if err != nil {
		return task, err
	}
	clickTask, err := s.client.GetTask(ticket.TicketID)
	if err != nil {
		return task, err
	}
	task = s.FormatTask(clickTask, ticket)
	comments, err := s.getAllTicketComments(ticketID)
	if err != nil {
		return task, err
	}
	task.Comments = comments
	return task, nil
}

func (s *Service) GetTicketsCount(orgID int32) (int32, error) {
	return s.dao.GetTicketsCount(orgID)
}

func (s *Service) EditTicket(user *models.User, ticketID string,
	req models.TicketUpdateRequest) error {
	ticket, err := s.dao.GetTicket(ticketID)
	if err != nil {
		return err
	}
	if ticket.OrgID != user.OrgId {
		return errors.New("invalid ticketID")
	}
	err = s.client.UpdateTask(ticketID, req)
	if len(req.Watcher) > 0 {
		err = s.dao.SetTicketWatcher(ticketID, req.Watcher)
		if err != nil {
			return err
		}
	}
	return err
}

func (s *Service) ReopenTask(user *models.User, ticketID string) error {
	ticket, err := s.dao.GetTicket(ticketID)
	if err != nil {
		return err
	}
	if ticket.OrgID != user.OrgId {
		return errors.New("invalid ticketID")
	}
	// TODO check if task can be open
	err = s.client.ReopenTask(ticketID)
	return err
}

func (s *Service) getTicketComments(ticketID, startID string) (
	[]models.ClickupTaskComment, error) {

	var result []models.ClickupTaskComment
	comments, err := s.client.GetTaskComments(ticketID, startID)
	for _, comment := range comments {
		parseResult, ok := models.ExtractComment(comment.CommentText)
		log.Debugf("raw content:%s parseResult:%+v ok:%v", comment.CommentText, parseResult, ok)
		if ok {
			comment.CommentText = parseResult.Content
			comment.User = models.ClickupCommentUser{
				ID:       int64(parseResult.UserID),
				Username: parseResult.Username,
				Email:    parseResult.Username,
			}
		} else {
			comment.User = models.ClickupCommentUser{
				ID:       0,
				Username: "system",
				Email:    "system@suger.io",
			}
		}
		log.Debugf("formatted comment:%+v", comment)
		result = append(result, comment)
	}
	if err != nil {
		return result, err
	}

	return result, err
}

func (s *Service) getAllTicketComments(ticketID string) (
	[]models.ClickupTaskComment, error) {
	var result []models.ClickupTaskComment

	startID := ""
	for {
		comments, err := s.getTicketComments(ticketID, startID)
		if err != nil {
			return result, err
		}
		if len(comments) < 1 {
			break
		}
		result = append(result, comments...)
		if len(comments) < TicketCommentPageSize {
			break
		} else {
			startID = comments[len(comments)-1].ID
		}
	}
	return result, nil
}

func (s *Service) GetTicketComments(user *models.User, ticketID, startID string) (
	[]models.ClickupTaskComment, error) {
	var result []models.ClickupTaskComment
	ticket, err := s.dao.GetTicket(ticketID)
	if err != nil {
		return result, err
	}
	if ticket.OrgID != user.OrgId {
		return result, errors.New("invalid ticketID")
	}
	return s.getTicketComments(ticketID, startID)
}

func (s *Service) CreateTicketComments(user *models.User, ticketID, commentText string) (
	string, error) {
	ticket, err := s.dao.GetTicket(ticketID)
	if err != nil {
		return "", err
	}
	if ticket.OrgID != user.OrgId {
		return "", errors.New("invalid ticketID")
	}
	parseResult := &models.ClickupCommentParseResult{
		Username: user.Email,
		UserID:   user.ID,
		Content:  commentText,
	}
	commentID, err := s.client.CreateTaskComment(ticketID, parseResult.Format())
	if err != nil {
		return "", err
	}
	return commentID, err
}

func (s *Service) SendNotification(user *models.User) {
	// TODO Send notification according to requirements

}
