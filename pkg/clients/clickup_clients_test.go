package clients

import (
	"testing"
)

func GetTestHandler() *clickupClient {
	return &clickupClient{
		AuthenticateKey: "pk_84243826_RMFURM9C270CLPWH47AEQA2SPJ7PTXJZ",
		ListId:          "901700864939",
	}
}

func TestCreateTask(t *testing.T) {
	h := GetTestHandler()
	task_id, err := h.CreateTask("hello", "world")
	if err != nil || task_id == "" {
		t.Errorf("fail to create task")
	}
}
