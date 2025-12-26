package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

type ErrorMsg struct {
	Field string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(s interface{}) []*ErrorMsg {
	var errors []*ErrorMsg

	err := Validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorMsg
			element.Field = err.Field()

			switch err.Tag() {
				case "required":
					element.Message = "Wajib diisi"
				case "email" :
					element.Message = "Format email salah"
				case "gt" :
					element.Message = fmt.Sprintf("Nilai harus lebih besar dari %s", err.Param())
				case "gte" :
					element.Message = fmt.Sprintf("Nilai harus lebih besar atau sama dengan %s", err.Param())
				case "min" :
					element.Message = fmt.Sprintf("Panjang minimal %s karakter", err.Param())
				case "max" :
					element.Message = fmt.Sprintf("Panjang maksimal %s karakter", err.Param())
				default:
					element.Message = "Tidak valid"
			}
			errors = append(errors, &element)
		}
	}
	return errors
}