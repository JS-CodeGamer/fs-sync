package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func GetValidator() *validator.Validate {
	return validate
}

// for _, err := range err.(validator.ValidationErrors) {
// 	fmt.Println(err.Namespace())
// 	fmt.Println(err.Field())
// 	fmt.Println(err.StructNamespace())
// 	fmt.Println(err.StructField())
// 	fmt.Println(err.Tag())
// 	fmt.Println(err.ActualTag())
// 	fmt.Println(err.Kind())
// 	fmt.Println(err.Type())
// 	fmt.Println(err.Value())
// 	fmt.Println(err.Param())
// 	fmt.Println()
// }
