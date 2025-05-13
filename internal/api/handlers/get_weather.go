package handlers

import (
	"errors"
	"net/http"

	"github.com/slbmax/ses-weather-app/internal/api/ctx"
	"github.com/slbmax/ses-weather-app/internal/api/requests"
	"github.com/slbmax/ses-weather-app/internal/api/responses"
	"github.com/slbmax/ses-weather-app/pkg/weatherapi"
	"gitlab.com/distributed_lab/ape"
)

func GetWeather(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewGetWeatherRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var (
		log           = ctx.GetLogger(r)
		weatherClient = ctx.GetWeatherClient(r)
	)

	weather, err := weatherClient.GetCurrentWeather(request.City)
	if err != nil {
		if errors.Is(err, weatherapi.ErrCityNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			log.WithError(err).Error("failed to get weather data")
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	response := responses.NewGetWeatherResponse(weather.CurrentWeather)
	ape.Render(w, response)
}
