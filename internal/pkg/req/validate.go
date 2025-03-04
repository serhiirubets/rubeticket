package req

import (
	"github.com/go-playground/validator/v10"
	"time"
)

func validateDateISO8601(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()

	// try parse to ISO 8601
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

func IsValid[T any](payload T) error {
	validate := validator.New()
	err := validate.RegisterValidation("date_iso8601", validateDateISO8601)
	if err != nil {
		return err
	}
	err = validate.Struct(payload)
	return err
}
