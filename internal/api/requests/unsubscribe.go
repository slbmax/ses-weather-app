package requests

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type UnsubscribeRequest struct {
	Token string
}

func (c *UnsubscribeRequest) Validate() error {
	return validation.Validate(c.Token, validation.Required, validation.Match(tokenRegex))
}

func NewUnsubscribeRequest(r *http.Request) (*UnsubscribeRequest, error) {
	request := &UnsubscribeRequest{Token: chi.URLParam(r, TokenParam)}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	return request, nil
}
