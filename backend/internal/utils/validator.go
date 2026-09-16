package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationErrorDetail struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

func ValidateStruct(s interface{}) []ValidationErrorDetail {
	var errs []ValidationErrorDetail
	if err := validate.Struct(s); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			message := fmt.Sprintf("Field '%s' is invalid on condition '%s'", field, err.Tag())

			switch err.Tag() {
			case "required":
				message = fmt.Sprintf("Bidang %s wajib diisi", field)
			case "email":
				message = fmt.Sprintf("Format email %s tidak valid", field)
			case "min":
				message = fmt.Sprintf("Panjang minimal %s adalah %s karakter/nilai", field, err.Param())
			case "max":
				message = fmt.Sprintf("Panjang maksimal %s adalah %s karakter/nilai", field, err.Param())
			case "gt":
				message = fmt.Sprintf("Nilai %s harus lebih besar dari %s", field, err.Param())
			case "gte":
				message = fmt.Sprintf("Nilai %s tidak boleh kurang dari %s", field, err.Param())
			case "oneof":
				message = fmt.Sprintf("Pilihan %s harus salah satu dari: %s", field, err.Param())
			}

			errs = append(errs, ValidationErrorDetail{
				Field:   field,
				Tag:     err.Tag(),
				Message: message,
			})
		}
	}
	return errs
}
