package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/slbmax/ses-weather-app/internal/api/requests"
	"github.com/slbmax/ses-weather-app/internal/api/responses"
	"github.com/slbmax/ses-weather-app/internal/database"
	subsMock "github.com/slbmax/ses-weather-app/internal/database/mock"
	mailerMock "github.com/slbmax/ses-weather-app/internal/mailer/mock"
	"github.com/slbmax/ses-weather-app/pkg/weatherapi"
	weatherApiMock "github.com/slbmax/ses-weather-app/pkg/weatherapi/mock"
	"github.com/stretchr/testify/mock"
	"gitlab.com/distributed_lab/logan/v3"
)

var (
	server           *httptest.Server
	subscriptionMock *subsMock.MockSubscriptionsQ
	weatherMock      *weatherApiMock.MockWeatherProvider
	mailMock         *mailerMock.MockMailer
)

func resetMocks() {
	subscriptionMock.Calls = []mock.Call{}
	subscriptionMock.Mock = mock.Mock{}

	weatherMock.Calls = []mock.Call{}
	weatherMock.Mock = mock.Mock{}

	mailMock.Calls = []mock.Call{}
	mailMock.Mock = mock.Mock{}
}

func TestMain(m *testing.M) {
	subscriptionMock = &subsMock.MockSubscriptionsQ{}
	weatherMock = &weatherApiMock.MockWeatherProvider{}
	mailMock = &mailerMock.MockMailer{}

	db := subsMock.NewDatabase(subscriptionMock)
	srv := NewServer(
		nil, // won't be even used
		weatherMock,
		db,
		mailMock,
		logan.New().Level(logan.ErrorLevel), // ignoring logging middleware
	)
	server = httptest.NewServer(srv.requestHandler())

	code := m.Run()

	server.Close()
	os.Exit(code)
}

func TestServer_Weather(t *testing.T) {
	stringPtr := func(s string) *string {
		return &s
	}

	testCases := map[string]struct {
		preparation    func()
		cleanup        func()
		city           *string
		expectedStatus int
		response       *responses.WeatherResponse
	}{
		"must 400 (missing url param)": {
			expectedStatus: http.StatusBadRequest,
		},
		"must 400 (empty city name)": {
			city:           stringPtr(""),
			expectedStatus: http.StatusBadRequest,
		},
		"must 404 (city not found)": {
			preparation: func() {
				weatherMock.On("GetCurrentWeather", "non-existent").Return(nil, weatherapi.ErrCityNotFound)
			},
			cleanup: func() {
				weatherMock.AssertExpectations(t)
				resetMocks()
			},
			city:           stringPtr("non-existent"),
			expectedStatus: http.StatusNotFound,
		},
		"must 500 (unknown error)": {
			preparation: func() {
				weatherMock.On("GetCurrentWeather", "London").Return(nil, errors.New("unknown error"))
			},
			cleanup: func() {
				weatherMock.AssertExpectations(t)
				resetMocks()
			},
			city:           stringPtr("London"),
			expectedStatus: http.StatusInternalServerError,
		},
		"must 200 (valid response)": {
			preparation: func() {
				weatherMock.On("GetCurrentWeather", "London").Return(&weatherapi.WeatherCurrentResponse{
					CurrentWeather: weatherapi.CurrentWeather{
						Temperature: 25,
						Humidity:    60,
						Condition: weatherapi.WeatherCondition{
							Text: "Sunny",
						},
					},
				}, nil)
			},
			cleanup: func() {
				weatherMock.AssertExpectations(t)
				resetMocks()
			},
			city:           stringPtr("London"),
			expectedStatus: http.StatusOK,
			response: &responses.WeatherResponse{
				Temperature: 25,
				Humidity:    60,
				Description: "Sunny",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.preparation != nil {
				tc.preparation()
			}

			endpoint := server.URL + "/api/weather"
			if tc.city != nil {
				endpoint += "?city=" + *tc.city
			}

			response, err := http.Get(endpoint)
			if tc.cleanup != nil {
				tc.cleanup()
			}
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}

			if response.StatusCode != tc.expectedStatus {
				t.Fatalf("expected status %d, got %d", tc.expectedStatus, response.StatusCode)
			}

			if tc.response != nil {
				var resp responses.WeatherResponse
				if err = json.NewDecoder(response.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if resp != *tc.response {
					t.Fatalf("expected response %+v, got %+v", *tc.response, resp)
				}
			}
		})
	}
}

func TestServer_Subscribe(t *testing.T) {
	testCases := map[string]struct {
		preparation    func()
		call           func() (*http.Response, error)
		expectedStatus int
		cleanup        func()
	}{
		"must 400 (missing required fields)": {
			call: func() (*http.Response, error) {
				return http.PostForm(server.URL+"/api/subscribe", nil)
			},
			expectedStatus: http.StatusBadRequest,
		},
		"must 400 (missing required fields in json body)": {
			call: func() (*http.Response, error) {
				return http.Post(server.URL+"/api/subscribe", "application/json", nil)
			},
			expectedStatus: http.StatusBadRequest,
		},
		"must 400 (invalid email)": {
			call: func() (*http.Response, error) {
				return http.PostForm(server.URL+"/api/subscribe", url.Values{
					"email":     {"invalid-email"},
					"city":      {"New York"},
					"frequency": {"daily"},
				})
			},
			expectedStatus: http.StatusBadRequest,
		},
		"must 409 (subscription already exists) ": {
			preparation: func() {
				subscriptionMock.On("Insert", mock.Anything).Return(int64(0), database.ErrSubscriptionExists)
			},
			call: func() (*http.Response, error) {
				return http.PostForm(server.URL+"/api/subscribe", url.Values{
					"email":     {"max@gmail.com"},
					"city":      {"New York"},
					"frequency": {"daily"},
				})
			},
			expectedStatus: http.StatusConflict,
			cleanup: func() {
				subscriptionMock.AssertExpectations(t)
				resetMocks()
			},
		},
		"must 404 (city not found error) ": {
			preparation: func() {
				subscriptionMock.On("Insert", mock.Anything).Return(int64(0), nil)
				weatherMock.On("GetCurrentWeather", "New York").Return(nil, weatherapi.ErrCityNotFound)
			},
			call: func() (*http.Response, error) {
				return http.PostForm(server.URL+"/api/subscribe", url.Values{
					"email":     {"max@gmail.com"},
					"city":      {"New York"},
					"frequency": {"daily"},
				})
			},
			expectedStatus: http.StatusNotFound,
			cleanup: func() {
				subscriptionMock.AssertExpectations(t)
				weatherMock.AssertExpectations(t)
				resetMocks()
			},
		},
		"must 500 (unknown error) ": {
			preparation: func() {
				subscriptionMock.On("Insert", mock.Anything).Return(int64(0), errors.New("db error"))
			},
			call: func() (*http.Response, error) {
				return http.PostForm(server.URL+"/api/subscribe", url.Values{
					"email":     {"max@gmail.com"},
					"city":      {"New York"},
					"frequency": {"daily"},
				})
			},
			expectedStatus: http.StatusInternalServerError,
			cleanup: func() {
				subscriptionMock.AssertExpectations(t)
				resetMocks()
			},
		},
		"must 500 (email sending error) ": {
			preparation: func() {
				subscriptionMock.On("Insert", mock.Anything).Return(int64(0), nil)
				weatherMock.On("GetCurrentWeather", "New York").Return(&weatherapi.WeatherCurrentResponse{}, nil)
				mailMock.On("SendConfirmationEmail", "max@gmail.com", mock.Anything).Return(errors.New("error"))

			},
			call: func() (*http.Response, error) {
				return http.PostForm(server.URL+"/api/subscribe", url.Values{
					"email":     {"max@gmail.com"},
					"city":      {"New York"},
					"frequency": {"daily"},
				})
			},
			expectedStatus: http.StatusInternalServerError,
			cleanup: func() {
				subscriptionMock.AssertExpectations(t)
				resetMocks()
			},
		},
		"must 200 (url val)": {
			preparation: func() {
				subscriptionMock.On("Insert", mock.Anything).Return(int64(1), nil)
				mailMock.On("SendConfirmationEmail", "max@gmail.com", mock.Anything).Return(nil)
				weatherMock.On("GetCurrentWeather", "New York").Return(&weatherapi.WeatherCurrentResponse{}, nil)
			},
			call: func() (*http.Response, error) {
				return http.PostForm(server.URL+"/api/subscribe", url.Values{
					"email":     {"max@gmail.com"},
					"city":      {"New York"},
					"frequency": {"daily"},
				})
			},
			expectedStatus: http.StatusOK,
			cleanup: func() {
				subscriptionMock.AssertExpectations(t)
				mailMock.AssertExpectations(t)
				weatherMock.AssertExpectations(t)
				resetMocks()
			},
		},
		"must 200 (json body)": {
			preparation: func() {
				subscriptionMock.On("Insert", mock.Anything).Return(int64(1), nil)
				mailMock.On("SendConfirmationEmail", "max@gmail.com", mock.Anything).Return(nil)
				weatherMock.On("GetCurrentWeather", "New York").Return(&weatherapi.WeatherCurrentResponse{}, nil)
			},
			call: func() (*http.Response, error) {
				req, _ := json.Marshal(requests.SubscribeRequest{
					Email:     "max@gmail.com",
					City:      "New York",
					Frequency: "daily",
				})

				return http.Post(server.URL+"/api/subscribe", "application/json", bytes.NewReader(req))
			},
			expectedStatus: http.StatusOK,
			cleanup: func() {
				subscriptionMock.AssertExpectations(t)
				mailMock.AssertExpectations(t)
				weatherMock.AssertExpectations(t)
				resetMocks()
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.preparation != nil {
				tc.preparation()
			}

			response, err := tc.call()
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}

			if tc.cleanup != nil {
				tc.cleanup()
			}

			if response.StatusCode != tc.expectedStatus {
				t.Fatalf("expected status %d, got %d", tc.expectedStatus, response.StatusCode)
			}
		})
	}
}

func TestServer_Confirm(t *testing.T) {
	validToken := "00000000000000000000000000000000"
	testCases := map[string]struct {
		preparation    func()
		cleanup        func()
		token          string
		expectedStatus int
	}{
		"must 400 (invalid token)": {
			expectedStatus: http.StatusBadRequest,
			token:          "awe",
		},
		"must 400 (subscription already confirmed)": {
			preparation: func() {
				subscriptionMock.On("GetByToken", validToken).Return(&database.Subscription{
					Confirmed: true,
				}, nil)
			},
			cleanup: func() {
				subscriptionMock.AssertExpectations(t)
				resetMocks()
			},
			token:          validToken,
			expectedStatus: http.StatusBadRequest,
		},
		"must 404 (token not found)": {
			preparation: func() {
				subscriptionMock.On("GetByToken", validToken).Return(nil, nil)
			},
			cleanup: func() {
				subscriptionMock.AssertExpectations(t)
				resetMocks()
			},
			token:          validToken,
			expectedStatus: http.StatusNotFound,
		},
		"must 500 (mail sending error)": {
			preparation: func() {
				subscriptionMock.On("GetByToken", validToken).Return(&database.Subscription{
					Id:        1,
					Confirmed: false,
					Email:     "max@gmail.com",
				}, nil)
				subscriptionMock.On("UpdateConfirmed", int64(1), mock.Anything).Return(nil)
				mailMock.On("SendConfirmationSuccessEmail", "max@gmail.com", mock.Anything).Return(errors.New("error"))
			},
			cleanup: func() {
				subscriptionMock.AssertExpectations(t)
				mailMock.AssertExpectations(t)
				resetMocks()
			},
			token:          validToken,
			expectedStatus: http.StatusInternalServerError,
		},
		"must 200": {
			preparation: func() {
				subscriptionMock.On("GetByToken", validToken).Return(&database.Subscription{
					Id:        1,
					Confirmed: false,
					Email:     "max@gmail.com",
				}, nil)
				subscriptionMock.On("UpdateConfirmed", int64(1), mock.Anything).Return(nil)
				mailMock.On("SendConfirmationSuccessEmail", "max@gmail.com", mock.Anything).Return(nil)
			},
			cleanup: func() {
				subscriptionMock.AssertExpectations(t)
				mailMock.AssertExpectations(t)
				resetMocks()
			},
			token:          validToken,
			expectedStatus: http.StatusOK,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.preparation != nil {
				tc.preparation()
			}

			response, err := http.Get(server.URL + "/api/confirm/" + tc.token)
			if tc.cleanup != nil {
				tc.cleanup()
			}
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}

			if response.StatusCode != tc.expectedStatus {
				t.Fatalf("expected status %d, got %d", tc.expectedStatus, response.StatusCode)
			}
		})
	}
}

func TestServer_Unsubscribe(t *testing.T) {
	validToken := "00000000000000000000000000000000"
	testCases := map[string]struct {
		preparation    func()
		cleanup        func()
		token          string
		expectedStatus int
	}{
		"must 400 (invalid token)": {
			expectedStatus: http.StatusBadRequest,
			token:          "awe",
		},
		"must 404 (token not found)": {
			preparation: func() {
				subscriptionMock.On("DeleteByToken", validToken).Return(database.ErrNoRowsAffected)
			},
			cleanup: func() {
				subscriptionMock.AssertExpectations(t)
				resetMocks()
			},
			token:          validToken,
			expectedStatus: http.StatusNotFound,
		},
		"must 500 (unknown error)": {
			preparation: func() {
				subscriptionMock.On("DeleteByToken", validToken).Return(errors.New("error"))
			},
			cleanup: func() {
				subscriptionMock.AssertExpectations(t)
				resetMocks()
			},
			token:          validToken,
			expectedStatus: http.StatusInternalServerError,
		},
		"must 200": {
			preparation: func() {
				subscriptionMock.On("DeleteByToken", validToken).Return(nil)
			},
			cleanup: func() {
				subscriptionMock.AssertExpectations(t)
				resetMocks()
			},
			token:          validToken,
			expectedStatus: http.StatusOK,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if tc.preparation != nil {
				tc.preparation()
			}

			response, err := http.Get(server.URL + "/api/unsubscribe/" + tc.token)
			if tc.cleanup != nil {
				tc.cleanup()
			}
			if err != nil {
				t.Fatalf("failed to make request: %v", err)
			}

			if response.StatusCode != tc.expectedStatus {
				t.Fatalf("expected status %d, got %d", tc.expectedStatus, response.StatusCode)
			}
		})
	}
}
