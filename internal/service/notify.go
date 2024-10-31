package service

import (
	"context"
	"fmt"

	"example.com/example/config"
)

func (svc *Services) SendNotification(ctx context.Context) error {
	c := config.Get()
	baseUrl := fmt.Sprintf("http://%s:%d", c.Host, c.Port)

	_, err := svc.Resty.R().
		EnableTrace().
		Put(baseUrl + "/api/users/notify")

	if err != nil {
		return fmt.Errorf("fail to send notification: %w", err)
	}

	return nil
}
