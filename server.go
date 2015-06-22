package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/shkh/lastfm-go/lastfm"
	"fmt"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

func main() {
	viper.SetConfigName("config")
	viper.SetEnvPrefix("elaborate_charts")
	viper.BindEnv("api_key")
	viper.BindEnv("api_secret")

	r := gin.Default()
	r.GET("/weekly_chart", weeklyChart)
	r.Run(":8080")
}

func weeklyChart(c *gin.Context) {
	username := strings.ToLower(c.Query("username"))
	chartType := strings.ToLower(c.Query("chartType"))
	fromDate, err1 := time.Parse(dateLayout, c.Query("fromDate"))
	toDate, err2 := time.Parse(dateLayout, c.Query("toDate"))

	if chartType != "artist" && chartType != "album" && chartType != "track" {
		respondWithError(c, 400, "Unknown chart type: " + chartType)
		return
	}
	if err1 != nil || err2 != nil {
		respondWithError(c, 400, `Date be presented in "%Y-%m-%d" format`)
	}

	err := viper.ReadInConfig()
	if err != nil {
		respondWithError(c, 500, "Unable to load config")
		return
	}

	apiKey := viper.GetString("api_key")
	apiSecret := viper.GetString("api_secret")
	api := lastfm.New(apiKey, apiSecret)
	result, err := api.User.GetWeeklyAlbumChart(lastfm.P{
		"user": username,
		"from": fromDate.Unix(),
		"to": toDate.Unix(),
	})
	if err != nil {
		respondWithError(c, 200, fmt.Sprintf("Failed to get chart for %s: %s", toDate.Format(dateLayout), err))
	}
	c.String(200, result.User)
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(400, map[string]string{
		"error": message,
	})
}
