package requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/slbmax/ses-weather-app/internal/database"
)

const (
	formParamEmail     = "email"
	formParamCity      = "city"
	formParamFrequency = "frequency"
)

var (
	RegexpEmail = regexp.MustCompile("^\\S+@\\S+\\.\\S+$") // basic one, without overkill
)

type SubscribeRequest struct {
	Email     string                         `json:"email"`
	City      string                         `json:"city"`
	Frequency database.SubscriptionFrequency `json:"frequency"`
}

func (req *SubscribeRequest) Validate() error {
	if req == nil {
		return fmt.Errorf("request is nil")
	}

	return validation.Errors{
		formParamEmail: validation.Validate(req.Email,
			validation.Required,
			validation.Match(RegexpEmail).Error("invalid email format"),
		),
		formParamCity: validation.Validate(req.City,
			validation.Required,
			validation.Length(1, 100).Error("invalid city name"),
		),
		formParamFrequency: validation.Validate(req.Frequency,
			validation.Required,
			validation.By(func(value interface{}) error {
				if f, ok := value.(database.SubscriptionFrequency); ok && f.Valid() {
					return nil
				}
				return fmt.Errorf("invalid frequency value: %v", value)
			}),
		),
	}.Filter()
}

func NewSubscribeRequest(r *http.Request) (*SubscribeRequest, error) {
	var req *SubscribeRequest

	// as noted in the spec, two content types should be supported
	switch r.Header.Get("Content-Type") {
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			return nil, fmt.Errorf("failed to parse form data: %w", err)
		}
		req = &SubscribeRequest{
			Email:     r.PostFormValue(formParamEmail),
			City:      r.PostFormValue(formParamCity),
			Frequency: database.SubscriptionFrequency(r.PostFormValue(formParamFrequency)),
		}
	case "application/json":
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return nil, fmt.Errorf("failed to decode json body: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported content type: %s", r.Header.Get("Content-Type"))
	}

	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate request: %w", err)
	}

	return req, nil
}
