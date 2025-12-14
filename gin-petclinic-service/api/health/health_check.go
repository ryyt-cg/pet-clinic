// Package health provides health endpoints & services
package health

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Check struct {
	Status string `json:"status" required:"true"`
}

func (hc Check) Validate() error {
	return validation.ValidateStruct(&hc,
		validation.Field(&hc.Status, validation.Required),
	)
}
