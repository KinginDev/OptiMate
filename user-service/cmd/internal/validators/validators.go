// Package validators
package validators

import "github.com/asaskevich/govalidator"

// CustomValidator is a custom validator
type CustomValidator struct{}

// Validate validates the input
// It returns an error if the validation fails
func (cv *CustomValidator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}
