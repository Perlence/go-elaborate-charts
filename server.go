package main

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/shkh/lastfm-go/lastfm"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

type (
	weeklyChartRequest struct {
		username string
		fromDate int64
		toDate   int64
	}

	elaborateError struct {
		code    int
		message string
	}

	PlayCounts map[string]int64

	WeeklyChartResponse struct {
		Chart  PlayCounts `json:"Chart"`
		ToDate int64      `json:"ToDate"`
	}
)

func main() {
	viper.SetConfigName("config")
	viper.SetEnvPrefix("elaborate_charts")
	viper.BindEnv("api_key")
	viper.BindEnv("api_secret")

	r := gin.Default()
	r.GET("/get_weekly_artist_chart", getWeeklyArtistChart)
	r.GET("/get_weekly_album_chart", getWeeklyAlbumChart)
	r.GET("/get_weekly_track_chart", getWeeklyTrackChart)
	r.GET("/get_info", getInfo)
	r.Run(":8080")
}

func getWeeklyArtistChart(c *gin.Context) {
	request, err := getWeeklyChartParams(c)
	if err != nil {
		respondWithError(c, err.(*elaborateError))
		return
	}
	api, err := getApi()
	if err != nil {
		respondWithError(c, err.(*elaborateError))
		return
	}
	result, err := api.User.GetWeeklyArtistChart(lastfm.P{
		"user": request.username,
		"from": request.fromDate,
		"to":   request.toDate,
	})
	if err != nil {
		respondWithError(c, newElaborateError(409, "Failed to get weekly artist chart: %s", err))
		return
	}
	chart := getPlayCounts(result)
	response := &WeeklyChartResponse{chart, request.toDate}
	c.JSON(200, structs.Map(response))
}

func getWeeklyAlbumChart(c *gin.Context) {
	request, err := getWeeklyChartParams(c)
	if err != nil {
		respondWithError(c, err.(*elaborateError))
		return
	}
	api, err := getApi()
	if err != nil {
		respondWithError(c, err.(*elaborateError))
		return
	}
	result, err := api.User.GetWeeklyAlbumChart(lastfm.P{
		"user": request.username,
		"from": request.fromDate,
		"to":   request.toDate,
	})
	if err != nil {
		respondWithError(c, newElaborateError(409, "Failed to get weekly album chart: %s", err))
		return
	}
	chart := getPlayCounts(result)
	response := &WeeklyChartResponse{chart, request.toDate}
	c.JSON(200, structs.Map(response))
}

func getWeeklyTrackChart(c *gin.Context) {
	request, err := getWeeklyChartParams(c)
	if err != nil {
		respondWithError(c, err.(*elaborateError))
		return
	}
	api, err := getApi()
	if err != nil {
		respondWithError(c, err.(*elaborateError))
		return
	}
	result, err := api.User.GetWeeklyTrackChart(lastfm.P{
		"user": request.username,
		"from": request.fromDate,
		"to":   request.toDate,
	})
	if err != nil {
		respondWithError(c, newElaborateError(409, "Failed to get weekly track chart: %s", err))
		return
	}
	chart := getPlayCounts(result)
	response := &WeeklyChartResponse{chart, request.toDate}
	c.JSON(200, structs.Map(response))
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

func getWeeklyChartParams(c *gin.Context) (*weeklyChartRequest, error) {
	username := strings.ToLower(c.Query("username"))
	rawFromDate := c.Query("fromDate")
	rawToDate := c.Query("toDate")
	if username == "" {
		return nil, newElaborateError(400, "Parameter 'username' is missing")
	}
	if rawFromDate == "" {
		return nil, newElaborateError(400, "Parameter 'fromDate' is missing")
	}
	if rawToDate == "" {
		return nil, newElaborateError(400, "Parameter 'toDate' is missing")
	}
	fromDate, err1 := strconv.ParseInt(rawFromDate, 10, 64)
	toDate, err2 := strconv.ParseInt(rawToDate, 10, 64)
	if err1 != nil || err2 != nil {
		return nil, newElaborateError(400, "Date must be presented in Unix format")
	}
	return &weeklyChartRequest{username, fromDate, toDate}, nil
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
	c.JSON(err.code, map[string]string{
		"error": err.message,
	})
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
	chart := make(PlayCounts)
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
