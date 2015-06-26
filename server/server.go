package main

import (
	"fmt"
	"github.com/Perlence/go-elaborate-charts/common"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/spf13/viper"
	"strconv"
	"strings"
	"time"
)

const AllowOrigin = "*"
const DateLayout = "2006-01-02"

type elaborateError struct {
	code    int
	message string
}

func main() {
	viper.SetConfigName("config")
	viper.SetEnvPrefix("elaborate_charts")
	viper.BindEnv("api_key")
	viper.BindEnv("api_secret")

	r := gin.Default()
	if gin.Mode() == gin.DebugMode {
		r.Use(CORSMiddleware())
	}
	r.GET("/get_weekly_chart", getWeeklyChart)
	r.GET("/get_info", getInfo)
	r.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", AllowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Token")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func getWeeklyChart(c *gin.Context) {
	request, err := newWeeklyChartParams(c)
	if err != nil {
		respondWithError(c, err.(*elaborateError))
		return
	}
	api, err := getApi()
	if err != nil {
		respondWithError(c, err.(*elaborateError))
		return
	}
	params := lastfm.P{
		"user": request.Username,
		"from": request.FromDate,
		"to":   request.ToDate,
	}
	var result interface{}
	switch request.ChartType {
	case "artist":
		result, err = api.User.GetWeeklyArtistChart(params)
	case "album":
		result, err = api.User.GetWeeklyAlbumChart(params)
	case "track":
		result, err = api.User.GetWeeklyTrackChart(params)
	}
	toDate := time.Unix(request.ToDate, 0).Format(DateLayout)
	if err != nil {
		respondWithError(c, newElaborateError(409, "Failed to get chart for week starting at %s", toDate))
		return
	}
	chart := getPlayCounts(result)
	response := &common.WeeklyChartResponse{chart, request.ToDate}
	c.JSON(200, response)
}

func getInfo(c *gin.Context) {
	username := strings.ToLower(c.Query("username"))
	api, err := getApi()
	if err != nil {
		respondWithError(c, err.(*elaborateError))
		return
	}
	result, err := api.User.GetInfo(lastfm.P{"user": username})
	if err != nil {
		respondWithError(c, newElaborateError(409, "Failed to get user info: %s", err))
		return
	}
	c.JSON(200, result)
}

func newWeeklyChartParams(c *gin.Context) (*common.WeeklyChartRequest, error) {
	username := strings.ToLower(c.Query("username"))
	chartType := strings.ToLower(c.Query("chart_type"))
	rawFromDate := c.Query("from_date")
	rawToDate := c.Query("to_date")
	if username == "" {
		return nil, newElaborateError(400, "Parameter 'username' is missing")
	}
	if chartType == "" {
		return nil, newElaborateError(400, "Parameter 'chart_type' is missing")
	}
	if rawFromDate == "" {
		return nil, newElaborateError(400, "Parameter 'from_date' is missing")
	}
	if rawToDate == "" {
		return nil, newElaborateError(400, "Parameter 'to_date' is missing")
	}
	if chartType != "artist" && chartType != "album" && chartType != "track" {
		return nil, newElaborateError(400, "Unrecognized chart type")
	}
	fromDate, err1 := strconv.ParseInt(rawFromDate, 10, 64)
	toDate, err2 := strconv.ParseInt(rawToDate, 10, 64)
	if err1 != nil || err2 != nil {
		return nil, newElaborateError(400, "Date must be presented in Unix format")
	}
	return &common.WeeklyChartRequest{username, chartType, fromDate, toDate}, nil
}

func getApi() (*lastfm.Api, error) {
	err := viper.ReadInConfig()
	if err != nil {
		return nil, newElaborateError(500, "Unable to load config")
	}

	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")
	return lastfm.New(apiKey, apiSecret), nil
}

func respondWithError(c *gin.Context, err *elaborateError) {
	c.JSON(err.code, &common.ErrorResponse{err.message})
}

func getPlayCounts(obj interface{}) map[string]int64 {
	entries := make([]map[string]interface{}, 0)
	switch obj.(type) {
	case lastfm.UserGetWeeklyArtistChart:
		for _, entry := range obj.(lastfm.UserGetWeeklyArtistChart).Artists {
			entries = append(entries, structs.Map(entry))
		}
	case lastfm.UserGetWeeklyAlbumChart:
		for _, entry := range obj.(lastfm.UserGetWeeklyAlbumChart).Albums {
			entries = append(entries, structs.Map(entry))
		}
	case lastfm.UserGetWeeklyTrackChart:
		for _, entry := range obj.(lastfm.UserGetWeeklyTrackChart).Tracks {
			entries = append(entries, structs.Map(entry))
		}
	}
	chart := make(common.PlayCounts)
	for _, entry := range entries {
		rawPlayCount := entry["PlayCount"].(string)
		name := entry["Name"].(string)
		playCount, err := strconv.ParseInt(rawPlayCount, 10, 64)
		if err == nil {
			chart[name] = playCount
		}
	}
	return chart
}

func newElaborateError(code int, message string, values ...interface{}) *elaborateError {
	return &elaborateError{code, fmt.Sprintf(message, values...)}
}

func (self *elaborateError) Error() string {
	return self.message
}
