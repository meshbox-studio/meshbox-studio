package activity

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/meshbox-studio/meshbox-studio/internal/auth"
)

func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "list-activity",
		Method:      http.MethodGet,
		Path:        "/api/activity",
		Summary:     "List activity feed",
		Description: "Returns the recent activity feed.",
		Tags:        []string{"Activity"},
		Middlewares: huma.Middlewares{auth.RequireSession(api)},
	}, listHandler)
}

type ListOutput struct {
	Body []Item
}

func listHandler(ctx context.Context, input *struct{}) (*ListOutput, error) {
	return &ListOutput{Body: store.List()}, nil
}