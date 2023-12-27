package clients

import (
	"fmt"
	"suger-clickup/pkg/models"

	"github.com/imroc/req/v3"
)

type ClickUpClient interface {
	CreateTask(string, string) (string, error)
	// TODO add rest methods
}

type clickupClient struct {
	AuthenticateKey string
	ListId          string
}

type taskCreateResponse struct {
	TaskId string `json:"id"`
}

func NewClickupHandler(key, listID string) ClickUpClient {
	return &clickupClient{
		AuthenticateKey: key,
		ListId:          listID,
	}
}

func (h *clickupClient) CreateTask(name, desp string) (string, error) {
	client := req.C().DevMode()
	var res taskCreateResponse
	task := models.Task{
		Name:        name,
		Description: desp,
	}
	resp, err := client.R().SetHeader("Authorization", h.AuthenticateKey).
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
