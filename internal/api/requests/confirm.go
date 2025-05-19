package requests

import (
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const TokenParam = "token"

var tokenRegex = regexp.MustCompile(`^[a-f0-9]{32}$`)

type ConfirmRequest struct {
	Token string
}

func (c *ConfirmRequest) Validate() error {
	return validation.Validate(c.Token, validation.Required, validation.Match(tokenRegex))
}

func NewConfirmRequest(r *http.Request) (*ConfirmRequest, error) {
	request := &ConfirmRequest{Token: chi.URLParam(r, TokenParam)}
	if err := request.Validate(); err != nil {
		return nil, err
	}

	return request, nil
}
