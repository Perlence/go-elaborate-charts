package common

type (
	PlayCounts map[string]int64

	WeeklyChartResponse struct {
		Chart  PlayCounts `json:"Chart"`
		ToDate int64      `json:"ToDate"`
	}
)
