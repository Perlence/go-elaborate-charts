package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Perlence/go-elaborate-charts/common"
	"github.com/gopherjs/jquery"
	"github.com/jinzhu/now"
	"github.com/rakyll/coop"
	"honnef.co/go/js/xhr"
	"html/template"
	"time"
)

var jQuery = jquery.NewJQuery

const (
	BackendUri = "http://127.0.0.1:8080"
)

func main() {
	app := NewApp()
	app.bindEvents()
}

type (
	App struct {
		alertTmpl               *template.Template
		formJQuery              jquery.JQuery
		usernameJQuery          jquery.JQuery
		chartTypeJQuery         jquery.JQuery
		numberOfPositionsJQuery jquery.JQuery
		timeframeJQuery         jquery.JQuery
		alertsJQuery            jquery.JQuery
	}

	Span struct {
		Start time.Time
		End   time.Time
	}
)

func NewApp() *App {
	settingsJQuery := jQuery("#settings-block")

	alertHtml := jQuery("#alert-template").Html()
	alertTmpl := template.Must(template.New("alert").Parse(alertHtml))

	formJQuery := settingsJQuery.Find("#form")

	usernameJQuery := formJQuery.Find("#username")
	chartTypeJQuery := formJQuery.Find("#chart-type")
	numberOfPositionsJQuery := formJQuery.Find("#number-of-positions")
	timeframeJQuery := formJQuery.Find("#timeframe")

	alertsJQuery := settingsJQuery.Find("#alerts")

	return &App{alertTmpl, formJQuery, usernameJQuery, chartTypeJQuery, numberOfPositionsJQuery, timeframeJQuery, alertsJQuery}
}

func (a *App) bindEvents() {
	a.formJQuery.On(jquery.SUBMIT, a.submit)
}

func (a *App) submit(e jquery.Event) bool {
	a.clearAlerts()

	uri := BackendUri + "/get_weekly_chart"
	queries, err := a.prepareRequests()
	if err != nil {
		a.showAlert("danger", "Bad request", err)
		return false
	}

	charts := make([]common.WeeklyChartResponse, 0)
	funcs := make([]func(), 0)
	for _, query := range queries {
		funcs = append(funcs, func(query common.WeeklyChartRequest) func() {
			return func() {
				req := xhr.NewRequest("GET", uri+"?"+query.Values().Encode())
				req.ResponseType = xhr.Text
				err := req.Send(nil)
				if err != nil {
					a.showAlert("danger", "Failed to get weekly charts", err)
					return
				}
				if 400 <= req.Status && req.Status <= 599 {
					var errorResponse common.ErrorResponse

					err = json.Unmarshal([]byte(req.ResponseText), &errorResponse)
					if err != nil {
						a.showAlert("danger", "Failed to parse weekly charts", err)
						return
					}
					a.showAlert("danger", "Server error", errors.New(errorResponse.Error))

					return
				}
				var chart common.WeeklyChartResponse

				err = json.Unmarshal([]byte(req.ResponseText), &chart)
				if err != nil {
					a.showAlert("danger", "Failed to parse weekly charts", err)
					return
				}

				charts = append(charts, chart)
			}
		}(query))
	}
	done := coop.All(funcs...)
	go func() {
		<-done
		// for _, chart := range charts {
		// 	fmt.Println(chart.ToDate)
		// }
	}()

	return false
}

func (a *App) prepareRequests() ([]common.WeeklyChartRequest, error) {
	username := a.usernameJQuery.Val()
	timeframe := a.timeframeJQuery.Val()
	chartType := a.chartTypeJQuery.Val()
	var fromDate, toDate time.Time
	toDate = time.Now().UTC()
	switch timeframe {
	case "last-7-days":
		fromDate = toDate.AddDate(0, 0, -7*2)
	case "last-month":
		fromDate = toDate.AddDate(0, -2, 0)
	case "last-3-months":
		fromDate = toDate.AddDate(0, -6, 0)
	case "last-6-months":
		fromDate = toDate.AddDate(0, -12, 0)
	case "last-12-months":
		fromDate = toDate.AddDate(0, -24, 0)
	case "overall":
		fromDate = time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC)
	default:
		return nil, errors.New("Unrecognized time frame: " + timeframe)
	}
	spans := dateSpanRange(fromDate, toDate, 0, 0, 7)
	queries := make([]common.WeeklyChartRequest, 0)
	for _, span := range spans {
		queries = append(
			queries,
			common.WeeklyChartRequest{
				username,
				chartType,
				span.Start.Unix(),
				span.End.Unix(),
			},
		)
	}
	return queries, nil
}

func (a *App) showAlert(style, reason string, err error) {
	alertData := struct {
		Style   string
		Reason  string
		Message string
	}{style, reason, err.Error()}

	var b bytes.Buffer
	a.alertTmpl.Execute(&b, alertData)
	strAlertTmpl := b.String()

	a.alertsJQuery.Append(strAlertTmpl)
}

func (a *App) clearAlerts() {
	a.alertsJQuery.Empty()
}

func dateSpanRange(start, end time.Time, years, months, days int) []Span {
	result := make([]Span, 0)
	var s, e time.Time
	s = now.New(start).BeginningOfWeek().Add(time.Duration(12) * time.Hour)
	for s.Before(end) {
		e = s.AddDate(years, months, days)
		if e.After(end) {
			e = end
		}
		result = append(result, Span{s, e})
		s = e
	}
	return result
}
