package httpsvc

import "github.com/go-playground/validator/v10"

type response struct {
	Status   string         `json:"status"` // success/failed
	Metadata map[string]any `json:"metadata"`
	Data     any            `json:"data"`
}

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}
