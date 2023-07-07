package internal

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var validators = map[string]func(validator.FieldLevel) bool{
	"validateEmail":    validateEmail,
	"validateUsername": validateUsername,
}

func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for key := range validators {
			err := v.RegisterValidation(key, validators[key])
			if err != nil {
				PrintFatal(fmt.Sprintf("error occurred while registering validator %s", key), err)
			}
		}
	}
}

func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if email == "" {
		return true
	}

	regex := regexp.MustCompile(`^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`)
	return regex.MatchString(email)
}

func validateUsername(fl validator.FieldLevel) bool {
	min := 3
	max := 50
	username := fl.Field().String()
	if username == "" {
		return true
	}

	return len(username) >= min && len(username) < max
}
