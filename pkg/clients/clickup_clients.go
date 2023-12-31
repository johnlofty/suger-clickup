package clients

import (
	"errors"
	"fmt"
	"suger-clickup/pkg/models"

	"github.com/gofiber/fiber/v2"
	"github.com/imroc/req/v3"
)

var ErrTaskNotFound = errors.New("task not found")

type ClickUpClient interface {
	CreateTask(task *models.Task) (string, error)
	GetTask(taskID string) (models.ClickupTask, error)
	ReopenTask(taskID string) error
	UpdateTask(taskID string, req models.TicketUpdateRequest) error

	GetList() (models.ClickupListResponse, error)
	GetTaskComments(taskID, startID string) ([]models.ClickupTaskComment, error)
	CreateTaskComment(taskID, commentText string) (string, error)
}

type clickupClient struct {
	AuthenticateKey string
	ListId          string
	client          *req.Client
}

type taskCreateResponse struct {
	TaskId string `json:"id"`
}

func NewClickupHandler(key, listID string) ClickUpClient {
	return &clickupClient{
		AuthenticateKey: key,
		ListId:          listID,
		client:          req.C().DevMode(),
	}
}

func (h *clickupClient) CreateTask(task *models.Task) (string, error) {
	var res taskCreateResponse
	resp, err := h.client.R().SetHeader("Authorization", h.AuthenticateKey).
		SetSuccessResult(&res).
		SetBody(task).
		Post(fmt.Sprintf("https://api.clickup.com/api/v2/list/%s/task", h.ListId))
	if err != nil {
		return "", err
	}
	if resp.IsErrorState() {
		return "", fmt.Errorf("request fail:%d ", resp.StatusCode)
	}
	return res.TaskId, nil
}

func (h *clickupClient) GetAccessibleCustomFields() {
	resp, err := h.client.R().SetHeader("Authorization", h.AuthenticateKey).
		Get(fmt.Sprintf("https://api.clickup.com/api/v2/list/%s/field", h.ListId))

	fmt.Println("resp:", resp, " err", err)
}

func (h *clickupClient) GetTask(taskID string) (models.ClickupTask, error) {
	var task models.ClickupTask
	resp, err := h.client.R().SetHeader("Authorization", h.AuthenticateKey).
		SetSuccessResult(&task).
		Get(fmt.Sprintf("https://api.clickup.com/api/v2/task/%s", taskID))

	if err != nil {
		return task, err
	}

	if resp.IsErrorState() {
		if resp.StatusCode == fiber.StatusNotFound {
			return task, ErrTaskNotFound
		}
		return task, fmt.Errorf("request fail:%d ", resp.StatusCode)
	}

	return task, nil
}

func (h *clickupClient) UpdateTask(taskID string,
	req models.TicketUpdateRequest) error {
	task := models.Task{
		Description: req.Description,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
	}
	resp, err := h.client.R().SetHeader("Authorization", h.AuthenticateKey).
		SetBody(task).
		Put(fmt.Sprintf("https://api.clickup.com/api/v2/task/%s", taskID))
	if err != nil {
		return err
	}
	if resp.IsErrorState() {
		return fmt.Errorf("request fail:%d ", resp.StatusCode)
	}
	return nil
}

func (h *clickupClient) UpdateTaskDescription(taskID, description string) error {
	task := models.Task{
		Description: description,
	}
	resp, err := h.client.R().SetHeader("Authorization", h.AuthenticateKey).
		SetBody(task).
		Put(fmt.Sprintf("https://api.clickup.com/api/v2/task/%s", taskID))
	if err != nil {
		return err
	}
	if resp.IsErrorState() {
		return fmt.Errorf("request fail:%d ", resp.StatusCode)
	}
	return nil
}

func (h *clickupClient) UpdateTaskDueDate(taskID string, dueDate int64) error {
	task := models.Task{
		DueDate: dueDate,
	}
	resp, err := h.client.R().SetHeader("Authorization", h.AuthenticateKey).
		SetBody(task).
		Put(fmt.Sprintf("https://api.clickup.com/api/v2/task/%s", taskID))
	if err != nil {
		return err
	}
	if resp.IsErrorState() {
		return fmt.Errorf("request fail:%d ", resp.StatusCode)
	}
	return nil
}

func (h *clickupClient) ReopenTask(taskID string) error {
	task := models.Task{
		Status: "Open",
	}
	resp, err := h.client.R().SetHeader("Authorization", h.AuthenticateKey).
		SetBody(task).
		Put(fmt.Sprintf("https://api.clickup.com/api/v2/task/%s", taskID))
	if err != nil {
		return err
	}
	if resp.IsErrorState() {
		return fmt.Errorf("request fail:%d ", resp.StatusCode)
	}
	return nil
}

func (h *clickupClient) GetSpace() (models.ClickupTask, error) {
	var task models.ClickupTask
	resp, err := h.client.R().SetHeader("Authorization", h.AuthenticateKey).
		SetSuccessResult(&task).
		Get(fmt.Sprintf("https://api.clickup.com/api/v2/task/%s", h.ListId))

	if err != nil {
		return task, err
	}

	if resp.IsErrorState() {
		return task, fmt.Errorf("request fail:%d ", resp.StatusCode)
	}

	return task, nil
}

func (h *clickupClient) GetList() (models.ClickupListResponse, error) {
	var task models.ClickupListResponse
	resp, err := h.client.R().SetHeader("Authorization", h.AuthenticateKey).
		SetSuccessResult(&task).
		Get(fmt.Sprintf("https://api.clickup.com/api/v2/list/%s", h.ListId))
	if err != nil {
		return task, err
	}

	if resp.IsErrorState() {
		return task, fmt.Errorf("request fail:%d ", resp.StatusCode)
	}

	return task, nil
}

func (h *clickupClient) GetTaskComments(taskID, startID string) (
	[]models.ClickupTaskComment, error) {
	var task models.ClickupTaskCommentsResponse
	params := map[string]string{
		"task_id":  taskID,
		"start_id": startID,
	}
	_, err := h.client.R().SetHeader("Authorization", h.AuthenticateKey).
		SetSuccessResult(&task).
		SetPathParams(params).
		Get("https://api.clickup.com/api/v2/task/{task_id}/comment?start_id={start_id}")
	if err != nil {
		return nil, err
	}
	return task.Comments, nil
}

func (h *clickupClient) CreateTaskComment(taskID, commentText string) (
	string, error) {
	var result models.ClickupCreateTaskCommentResponse
	params := map[string]string{
		"task_id": taskID,
	}
	req := models.ClickupCreateTaskCommentRequest{
		CommentText: commentText,
	}
	_, err := h.client.R().SetHeader("Authorization", h.AuthenticateKey).
		SetSuccessResult(&result).
		SetPathParams(params).
		SetBody(req).
		Post("https://api.clickup.com/api/v2/task/{task_id}/comment")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", result.ID), nil
}
