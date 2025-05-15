package weatherapi

type WeatherCondition struct {
	Text string `json:"text"`
}

type Location struct {
	Name string `json:"name"`
}

type CurrentWeather struct {
	Temperature float32          `json:"temp_c"`
	Humidity    uint8            `json:"humidity"`
	Condition   WeatherCondition `json:"condition"`
}

type WeatherCurrentResponse struct {
	// adding only necessary fields for brevity
	CurrentWeather CurrentWeather `json:"current"`
	Location       Location       `json:"location"`
}
