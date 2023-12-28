package clients

import (
	"fmt"
	"suger-clickup/pkg/models"

	"github.com/imroc/req/v3"
)

type ClickUpClient interface {
	CreateTask(task *models.Task) (string, error)
	GetTask(taskID string) (models.ClickupTask, error)
	// TODO add rest methods
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
		return task, fmt.Errorf("request fail:%d ", resp.StatusCode)
	}

	return task, nil
}
