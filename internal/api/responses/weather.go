package responses

import "github.com/slbmax/ses-weather-app/pkg/weatherapi"

type GetWeatherResponse struct {
	Temperature float32 `json:"temperature"`
	Humidity    uint8   `json:"humidity"`
	Description string  `json:"description"`
}

func NewGetWeatherResponse(weather weatherapi.CurrentWeather) GetWeatherResponse {
	return GetWeatherResponse{
		Temperature: weather.Temperature,
		Humidity:    weather.Humidity,
		Description: weather.Condition.Text,
	}
}
