package clients

import (
	"fmt"
	"suger-clickup/pkg/models"
	"suger-clickup/pkg/settings"
)

var statusMap map[string]models.ClickupTaskStatus

func Setup() {
	setupClickStatus()
}

func setupClickStatus() {
	conf := settings.Get()
	client := NewClickupHandler(
		conf.ClickupConfig.Secret,
		conf.ClickupConfig.ListId,
	)
	listInfo, err := client.GetList()
	if err != nil {
		panic(fmt.Sprintf("fail to get status:%v", err))
	}

	statusMap = make(map[string]models.ClickupTaskStatus)
	for _, status := range listInfo.Statuses {
		statusMap[status.Status] = status
	}
}
