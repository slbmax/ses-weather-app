package requests

import (
	"fmt"
	"net/http"
)

const queryParamCity = "city"

type WeatherRequest struct {
	City string
}

func NewWeatherRequest(r *http.Request) (*WeatherRequest, error) {
	query := r.URL.Query()
	city := query.Get(queryParamCity)
	if city == "" {
		return nil, fmt.Errorf("%q parameter is required", queryParamCity)
	}

	return &WeatherRequest{City: city}, nil
}
