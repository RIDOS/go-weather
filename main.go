package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

type GeoResponse struct {
	Results []LatLong `json:"results"`
}

type LatLong struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type WeatherResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
	Hourly    struct {
		Time          []string  `json:"time"`
		Temperature2m []float64 `json:"temperature_2m"`
	} `json:"hourly"`
}

type WeatherDisplay struct {
	City      string
	Forecasts []Forecast
}

type Forecast struct {
	Date        string
	Temperature string
}

/*
*

	Вывод координат города в формате JSON.
	@return string, error
	Пример:
	{
		"results": [
		{
			"id": 555813,
			"name": "Gora Isonkiventunturi",
			"latitude": 69.49955,
			"longitude": 30.96885,
			"elevation": 298.0,
			"feature_code": "MT",
			"country_code": "RU",
			"admin1_id": 524304,
			"timezone": "Europe/Moscow",
			"country_id": 2017370,
			"country": "Russia",
			"admin1": "Murmansk"
		}],
		"generationtime_ms": 0.5930662
	}
*/
func getLatLong(city string) (*LatLong, error) {
	endpoint := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=en&format=json", url.QueryEscape(city))
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error making request to Geo API: %w", err)
	}
	defer resp.Body.Close()

	var response GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	if len(response.Results) < 1 {
		return nil, errors.New("no results found")
	}

	return &response.Results[0], nil
}

/*
*

		Вывод погоды по координатам города в формате JSON.
		@return string, error

		Пример:
		{
		"latitude": 53.625,
		"longitude": 55.9375,
		"generationtime_ms": 0.06401538848876953,
		"utc_offset_seconds": 0,
		"timezone": "GMT",
		"timezone_abbreviation": "GMT",
		"elevation": 168.0,
		"hourly_units":
		{
			"time": "iso8601",
			"temperature_2m": "°C"
		},
		"hourly":
		{
			"time": ["2023-10-01T00:00", "2023-10-01T01:00", "2023-10-01T02:00", "2023-10-01T03:00", ...
			"temperature_2m": [6.6, 6.2, 5.6, 5.2, 6.0, 8.7, 11.4, 13.4, 14.8, 15.9, 16.5, 16.7, 16.5, 15.7, 14.0, ...
		}
	}
*/
func getWeather(latLong LatLong) (string, error) {
	endpoint := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.6f&longitude=%.6f&hourly=temperature_2m", latLong.Latitude, latLong.Longitude)
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("error making request to Weather API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

/*
*

	Вывод JSON: информации о погоде.
	@return JSON
*/
func getWeatherByQuery(c *gin.Context) {
	city := c.Query("city")

	latLong, err := getLatLong(city)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	weather, err := getWeather(*latLong)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"weather": weather})
}

/*
*

	Обработка и вывод в понятном формате.
	@return WeatherDisplay
*/
func extractWeatherData(city string, rawWeather string) (WeatherDisplay, error) {
	var weatherResponse WeatherResponse
	if err := json.Unmarshal([]byte(rawWeather), &weatherResponse); err != nil {
		return WeatherDisplay{}, fmt.Errorf("error decoding weather response: %w", err)
	}

	var forecasts []Forecast
	for i, t := range weatherResponse.Hourly.Time {
		date, err := time.Parse("2006-01-02T15:04", t)
		if err != nil {
			return WeatherDisplay{}, err
		}
		forecast := Forecast{
			Date:        date.Format("Mon 15:04"),
			Temperature: fmt.Sprintf("%.1f°C", weatherResponse.Hourly.Temperature2m[i]),
		}
		forecasts = append(forecasts, forecast)
	}
	return WeatherDisplay{
		City:      city,
		Forecasts: forecasts,
	}, nil
}

func main() {
	// Подключение шаблонов.
	r := gin.Default()
	r.LoadHTMLGlob("views/*")

	// Тест для локальной разработки самого сервиса.
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/weather", func(c *gin.Context) {
		city := c.Query("city")
		latlong, err := getLatLong(city)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		weather, err := getWeather(*latlong)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		weatherDisplay, err := extractWeatherData(city, weather)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "weather.html", weatherDisplay)
	})

	// Отдельный запрос с выводом в формтае JSON.
	r.GET("/weatherJSON", func(c *gin.Context) {
		city := c.Query("city")
		latlong, err := getLatLong(city)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		weather, err := getWeather(*latlong)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		weatherDisplay, err := extractWeatherData(city, weather)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"OK": weatherDisplay})
	})

	// Используемый порт.
	http.ListenAndServe(":3000", r)
}
