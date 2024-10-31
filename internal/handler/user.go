package handler

import (
	"context"
	"log/slog"
	"net/http"

	"example.com/example/internal/service"
	"example.com/example/lib/logging"
	"github.com/danielgtaylor/huma/v2"
)

type UserResponseBody struct {
	Name   string   `json:"name"`
	Emails []string `json:"emails"`
}

type UserCreateBody struct {
	Name   string   `json:"name" example:"john" required:"true"`
	Emails []string `json:"emails" example:"a@a.com,b@b.com" required:"true"`
}

type CreateUserInput struct {
	Body UserCreateBody `required:"true"`
}

type ListUserOutput struct {
	Body   []UserResponseBody
	Status int
}

type CreateUserOutput struct {
	Body   UserResponseBody
	Status int
}

func registerUser(api huma.API, svc service.AllServices) {
	huma.Register(api,
		huma.Operation{
			OperationID: "get-users",
			Method:      http.MethodGet,
			Path:        "/api/users",
			Summary:     "Get a bunch of users",
		}, func(ctx context.Context, input *struct{}) (*ListUserOutput, error) {
			users, err := svc.UserList(ctx)

			if err != nil {
				logging.Error("Failed to list story", slog.Any("error", err))
				return nil, huma.Error400BadRequest("fail")
			}

			var data []UserResponseBody
			for _, user := range users {
				var r UserResponseBody
				r.Name = user.Name
				r.Emails = user.Emails
				data = append(data, r)
			}

			status := http.StatusOK
			if len(data) == 0 {
				status = http.StatusNoContent
			}

			return &ListUserOutput{
				Body:   data,
				Status: status,
			}, nil
		},
	)

	huma.Register(api,
		huma.Operation{
			OperationID: "create-user",
			Method:      http.MethodPost,
			Path:        "/api/users",
			Summary:     "Create a user",
		}, func(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
			user, err := svc.UserCreate(ctx, input.Body.Name, input.Body.Emails)
			if err != nil {
				logging.Error("Failed to create user", slog.Any("error", err))
				return nil, huma.Error400BadRequest("fail")
			}

			var r UserResponseBody
			r.Name = user.Name
			r.Emails = user.Emails

			_ = svc.SendNotification(ctx)

			return &CreateUserOutput{
				Body:   r,
				Status: http.StatusCreated,
			}, nil
		},
	)

	huma.Register(api,
		huma.Operation{
			OperationID: "notify-user",
			Method:      http.MethodPut,
			Path:        "/api/users/notify",
			Summary:     "Notify a user",
		}, func(ctx context.Context, input *struct{}) (*struct{}, error) {
			logging.Info("send notification")
			return nil, nil
		},
	)
}
