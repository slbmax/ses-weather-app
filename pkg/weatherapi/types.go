package weatherapi

type CurrentWeather struct {
	Temperature float64 `json:"temp_c"`
	Humidity    int8    `json:"humidity"`
}

type WeatherCurrentResponse struct {
	// adding only necessary fields for brevity
	CurrentWeather CurrentWeather `json:"current"`
}
