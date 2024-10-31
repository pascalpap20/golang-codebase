package handler

import (
	"context"
	"log/slog"
	"net/http"

	"example.com/example/lib/logging"
	"git.govtechindonesia.id/inadigital/inatrace"
	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

func (h *Handler) GetUsers(ctx context.Context, input *struct{}) (*ListUserOutput, error) {
	users, err := h.svc.UserList(ctx)

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
}

func (h *Handler) CreateUser(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
	user, err := h.svc.UserCreate(ctx, input.Body.Name, input.Body.Emails)
	if err != nil {
		logging.Error("Failed to create user", slog.Any("error", err))
		return nil, huma.Error400BadRequest("fail")
	}

	var r UserResponseBody
	r.Name = user.Name
	r.Emails = user.Emails

	_ = h.svc.SendNotification(ctx)

	return &CreateUserOutput{
		Body:   r,
		Status: http.StatusCreated,
	}, nil
}

func (h *Handler) NotifyUser(ctx context.Context, input *struct{}) (*struct {
	Body []byte
}, error) {
	_, span := inatrace.Start(ctx, "sendNotification", trace.WithAttributes(attribute.String("id", "id")))
	defer span.End()

	return &struct {
		Body []byte
	}{
		Body: []byte("."),
	}, nil
}

func (h *Handler) RegisterUser(api huma.API) {
	huma.Register(api,
		huma.Operation{
			OperationID: "get-users",
			Method:      http.MethodGet,
			Path:        "/api/users",
			Summary:     "Get a bunch of users",
		}, h.GetUsers,
	)

	huma.Register(api,
		huma.Operation{
			OperationID: "create-user",
			Method:      http.MethodPost,
			Path:        "/api/users",
			Summary:     "Create a user",
		}, h.CreateUser,
	)

	huma.Register(api,
		huma.Operation{
			OperationID: "notify-user",
			Method:      http.MethodPut,
			Path:        "/api/users/notify",
			Summary:     "Notify a user",
		}, h.NotifyUser,
	)
}
