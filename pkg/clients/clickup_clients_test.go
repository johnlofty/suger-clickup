package clients

import (
	"suger-clickup/pkg/models"
	"testing"

	"github.com/imroc/req/v3"
)

func GetTestHandler() *clickupClient {
	return &clickupClient{
		AuthenticateKey: "pk_84243826_RMFURM9C270CLPWH47AEQA2SPJ7PTXJZ",
		ListId:          "901700864939",
		client:          req.C().DevMode(),
	}
}

func TestCreateTask(t *testing.T) {
	h := GetTestHandler()
	task := &models.Task{
		Name:        "hello",
		Description: "world",
	}
	task.AddOrgId("4")
	task.AddOrgName("storm")
	task.AddCreator("creator")
	task_id, err := h.CreateTask(task)
	if err != nil || task_id == "" {
		t.Errorf("fail to create task")
	}
}

func TestGetCustomeFields(t *testing.T) {
	h := GetTestHandler()
	h.GetAccessibleCustomFields()
	t.Errorf("fail ")
}

func TestGetTask(t *testing.T) {
	h := GetTestHandler()
	taskID := "86dqzahm0"
	task, err := h.GetTask(taskID)
	t.Errorf("getting task:%+v err:%v", task, err)
}
