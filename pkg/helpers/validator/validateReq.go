package validator

import (
	logger "github.com/CodeChefVIT/cookoff-10.0-be/pkg/logging"
	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func InitValidator() {
	Validate = validator.New()
}

func ValidatePayload(v any) error {
	if err := Validate.Struct(v); err != nil {
		logger.Warnf("%s", err.Error())
		return err
	}
	return nil
}
