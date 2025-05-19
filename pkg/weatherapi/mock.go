package weatherapi

type MockWeatherProvider struct{}

func NewMockWeatherProvider() WeatherProvider {
	return &MockWeatherProvider{}
}

func (m *MockWeatherProvider) GetCurrentWeather(city string) (*WeatherCurrentResponse, error) {
	if city == "" {
		return nil, ErrCityNotFound
	}

	return &WeatherCurrentResponse{
		Location: Location{
			Name: city,
		},
		CurrentWeather: CurrentWeather{
			Temperature: 20.1,
			Humidity:    60,
			Condition: WeatherCondition{
				Text: "Sunny",
			},
		},
	}, nil
}
