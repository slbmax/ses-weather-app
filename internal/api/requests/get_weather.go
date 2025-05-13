package requests

import (
	"fmt"
	"net/http"
)

const queryParamCity = "city"

type GetWeatherRequest struct {
	City string
}

func NewGetWeatherRequest(r *http.Request) (*GetWeatherRequest, error) {
	query := r.URL.Query()
	city := query.Get(queryParamCity)
	if city == "" {
		return nil, fmt.Errorf("%q parameter is required", queryParamCity)
	}

	return &GetWeatherRequest{City: city}, nil
}
