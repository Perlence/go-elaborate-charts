package common

import (
	"net/url"
	"strconv"
)

type (
	WeeklyChartRequest struct {
		Username  string `url:"username"`
		ChartType string `url:"chart_type"`
		FromDate  int64  `url:"from_date"`
		ToDate    int64  `url:"to_date"`
	}

	ErrorResponse struct {
		Error string
	}

	PlayCounts map[string]int64

	WeeklyChartResponse struct {
		Chart  PlayCounts `json:"Chart"`
		ToDate int64      `json:"ToDate"`
	}

	UserInfoRequest struct {
		Username string `url:"username"`
	}
)

func (r *WeeklyChartRequest) Values() url.Values {
	return url.Values{
		"username":   []string{r.Username},
		"chart_type": []string{r.ChartType},
		"from_date":  []string{strconv.FormatInt(r.FromDate, 10)},
		"to_date":    []string{strconv.FormatInt(r.ToDate, 10)},
	}
}
