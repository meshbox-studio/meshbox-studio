package stats

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-stats",
		Method:      "GET",
		Path:        "/api/stats",
		Summary:     "Get sidebar stats",
		Description: "Returns system information for the sidebar footer.",
		Tags:        []string{"Stats"},
	}, func(ctx context.Context, input *struct{}) (*StatsOutput, error) {
		return GetStats(), nil
	})
}

type StatsOutput struct {
	Body SidebarStats
}

func GetStats() *StatsOutput {
	s := GetSidebarStats()
	return &StatsOutput{Body: s}
}