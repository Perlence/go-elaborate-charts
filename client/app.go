package main

import (
	"bytes"
	"errors"
	"github.com/Perlence/go-elaborate-charts/common"
	"github.com/franela/goreq"
	"github.com/gopherjs/jquery"
	"github.com/rakyll/coop"
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
	uri, err := a.getRequestUri()
	if err != nil {
		a.showAlert("danger", "Bad request", err)
		return false
	}
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
				resp, err := goreq.Request{
					Uri:         uri,
					QueryString: query.Values(),
				}.Do()
				if err != nil {
					a.showAlert("danger", "Failed to get weekly charts", err)
					return
				}
				var chart common.WeeklyChartResponse
				resp.Body.FromJsonTo(&chart)
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

func (a *App) getRequestUri() (string, error) {
	chartType := a.chartTypeJQuery.Val()
	var uriPart string
	switch chartType {
	case "artist":
		uriPart = "/get_weekly_artist_chart"
	case "album":
		uriPart = "/get_weekly_album_chart"
	case "track":
		uriPart = "/get_weekly_track_chart"
	default:
		return "", errors.New("Unrecognized chart type: " + chartType)
	}
	return BackendUri + uriPart, nil
}

func (a *App) prepareRequests() ([]common.WeeklyChartRequest, error) {
	username := a.usernameJQuery.Val()
	timeframe := a.timeframeJQuery.Val()
	var fromDate, toDate time.Time
	toDate = time.Now()
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

	a.alertsJQuery.SetHtml(strAlertTmpl)
}

func dateSpanRange(start, end time.Time, years, months, days int) []Span {
	result := make([]Span, 0)
	var s, e time.Time
	s = start
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
