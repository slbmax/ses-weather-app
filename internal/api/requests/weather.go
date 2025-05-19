package requests

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const queryParamCity = "city"

type WeatherRequest struct {
	City string
}

func (r *WeatherRequest) Validate() error {
	return validation.Validate(r.City, validation.Required, validation.Length(1, 100).Error("invalid city name"))
}

func NewWeatherRequest(r *http.Request) (*WeatherRequest, error) {
	query := r.URL.Query()
	req := &WeatherRequest{query.Get(queryParamCity)}
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	return req, nil
}
