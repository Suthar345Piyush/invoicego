// Validator code here

package util

import "github.com/go-playground/validator"

var validate = validator.New()

// function for validating struct

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
