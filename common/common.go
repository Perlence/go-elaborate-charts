package common

import (
	"net/url"
	"strconv"
)

type (
	WeeklyChartRequest struct {
		Username string `url:"username"`
		FromDate int64  `url:"fromDate"`
		ToDate   int64  `url:"toDate"`
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
		"username": []string{r.Username},
		"fromDate": []string{strconv.FormatInt(r.FromDate, 10)},
		"toDate":   []string{strconv.FormatInt(r.ToDate, 10)},
	}
}
