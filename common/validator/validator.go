package validator

import (
	"log"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/SawitProRecruitment/UserService/common/errors"
)

const (
	ContentTypeJSON   = "json"
	ContentTypeStruct = "struct"
)

func Validate(content string, s interface{}) error {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	// Register the custom validation functions
	validate.RegisterValidation("phone", validatePhoneNumber)
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add("password", "{0} must contains uppercase, numeric, and special character minimal 1", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())

		return t
	})
	validate.RegisterTranslation("phone", trans, func(ut ut.Translator) error {
		return ut.Add("phone", "{0} must be start with +62, min 10 characters, max 13 characters", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field())

		return t
	})

	en_translations.RegisterDefaultTranslations(validate, trans)

	if content == ContentTypeJSON {
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get(content), ",", 2)[0]

			if name == "-" {
				return ""
			}

			return name
		})
	}

	// validate
	err := validate.Struct(s)
	if err != nil {
		switch ve := err.(type) {
		case validator.ValidationErrors:
			e := errors.BadRequest("")
			for _, err := range ve {
				et := err.Translate(trans)
				var msg = et
				if err.Field() != "" {
					log.Println(et)
					msg = strings.ReplaceAll(et, err.Field()+" ", "")
				}

				e.AddField(err.Field(), msg)
			}

			return e
		case *validator.InvalidValidationError:
			return ve
		}

		return err
	}

	return nil
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()
	// Define a regular expression pattern for an Indonesian phone number
	pattern := `^\+62\d{8,11}$`
	match, _ := regexp.MatchString(pattern, phoneNumber)
	return match
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 6 || len(password) > 64 {
		return false
	}

	// Check for at least one uppercase letter
	hasUppercase := false
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUppercase = true
			break
		}
	}

	// Check for at least one digit (number)
	hasDigit := false
	for _, char := range password {
		if unicode.IsDigit(char) {
			hasDigit = true
			break
		}
	}

	// Check for at least one special character
	hasSpecialChar := false
	specialChars := "!@#$%^&*()-_=+\\|[{]};:'\",<.>/?"
	for _, char := range password {
		if strings.ContainsRune(specialChars, char) {
			hasSpecialChar = true
			break
		}
	}

	return hasUppercase && hasDigit && hasSpecialChar
}
