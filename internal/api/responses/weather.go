package responses

import "github.com/slbmax/ses-weather-app/pkg/weatherapi"

type WeatherResponse struct {
	Temperature float32 `json:"temperature"`
	Humidity    uint8   `json:"humidity"`
	Description string  `json:"description"`
}

func NewWeatherResponse(weather weatherapi.CurrentWeather) WeatherResponse {
	return WeatherResponse{
		Temperature: weather.Temperature,
		Humidity:    weather.Humidity,
		Description: weather.Condition.Text,
	}
}
