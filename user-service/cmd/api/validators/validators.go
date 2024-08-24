package validators

import "github.com/asaskevich/govalidator"

type CustomValidator struct{}

func (cv *CustomValidator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}
