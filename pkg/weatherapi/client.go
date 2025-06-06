package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const baseUrl = "https://api.weatherapi.com/v1"

type WeatherProvider interface {
	GetCurrentWeather(city string) (*WeatherCurrentResponse, error)
}

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey}
}

func (c *Client) GetCurrentWeather(city string) (*WeatherCurrentResponse, error) {
	url := baseUrl + "/current.json?key=" + c.apiKey + "&q=" + city
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get current weather: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusBadRequest:
			// assuming 400 means city not found
			return nil, ErrCityNotFound
		default:
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	}

	var weatherResponse WeatherCurrentResponse
	if err = json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return nil, err
	}

	// weather api can return a different city name than requested (especially when auto-completing something)
	if strings.ToLower(weatherResponse.Location.Name) != strings.ToLower(city) {
		return nil, ErrCityNotFound
	}

	return &weatherResponse, nil
}
