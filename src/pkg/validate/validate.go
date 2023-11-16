package validate

import (
	"errors"
	"regexp"
	"slices"
	"strings"
	"time"

	"crdx.org/lighthouse/pkg/duration"
	"crdx.org/lighthouse/util/reflectutil"
	"crdx.org/lighthouse/util/stringutil"
	enLocale "github.com/go-playground/locales/en"
	universalTranslator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/samber/lo"
)

// Field represents a submitted HTML field. Error is the error (if any) and Value is the original
// value that was submitted in the request.
type Field struct {
	Error string
	Value string
	Name  string
}

type ValidatorMap map[string]func(string) error

var (
	validate   *validator.Validate
	translator universalTranslator.Translator
)

var ErrValidationFailed = errors.New("validation failed")

func init() {
	translator, _ = universalTranslator.New(enLocale.New()).GetTranslator("en")
	validate = validator.New(validator.WithRequiredStructEnabled())

	lo.Must0(enTranslations.RegisterDefaultTranslations(validate, translator))

	Register("timezone", "must be a valid timezone", func(value string) bool {
		_, err := time.LoadLocation(value)
		return err == nil
	})

	Register("mailaddr", `must be in the format "name <email>"`, func(value string) bool {
		return regexp.MustCompile("^.* <.*>$").Match([]byte(value))
	})

	Register("role", "must be a valid role", func(value string) bool {
		return slices.Contains([]string{"1", "2", "3"}, value)
	})

	Register("icon", "must be a valid icon", func(value string) bool {
		return regexp.MustCompile("^(duotone|solid|brands):.+$").Match([]byte(value))
	})

	Register("duration", "must be a valid duration", duration.Valid)

	RegisterWithParam("dmin", "must be at least {0}", func(value string, min string) bool {
		minDuration, ok := duration.Parse(min)
		return ok && lo.Must(duration.Parse(value)) >= minDuration
	})

	RegisterWithParam("dmax", "must be at most {0}", func(value string, max string) bool {
		maxDuration, ok := duration.Parse(max)
		return ok && lo.Must(duration.Parse(value)) <= maxDuration
	})
}

// Register registers a custom validation function.
func Register(name string, err string, f func(string) bool) {
	lo.Must0(validate.RegisterValidation(name, func(field validator.FieldLevel) bool {
		return f(field.Field().String())
	}, false))

	registerFn := func(translator universalTranslator.Translator) error {
		return translator.Add(name, err, true)
	}

	translateFn := func(translator universalTranslator.Translator, _ validator.FieldError) string {
		return lo.Must(translator.T(name))
	}

	lo.Must0(validate.RegisterTranslation(name, translator, registerFn, translateFn))
}

// RegisterWithParam registers a custom validation function that takes a parameter.
func RegisterWithParam(name string, err string, f func(value string, param string) bool) {
	lo.Must0(validate.RegisterValidation(name, func(field validator.FieldLevel) bool {
		return f(field.Field().String(), field.Param())
	}, false))

	registerFn := func(translator universalTranslator.Translator) error {
		return translator.Add(name, err, true)
	}

	translateFn := func(translator universalTranslator.Translator, e validator.FieldError) string {
		return lo.Must(translator.T(name, e.Param()))
	}

	lo.Must0(validate.RegisterTranslation(name, translator, registerFn, translateFn))
}

// Fields returns the initial data needed by the template to render the fields: the field names.
func Fields[T any]() map[string]Field {
	fields := map[string]Field{}

	structValue := reflectutil.GetValue(new(T))

	for i := 0; i < structValue.NumField(); i++ {
		fieldName := structValue.Type().Field(i).Name
		tagValue := structValue.Type().Field(i).Tag.Get("form")

		fields[fieldName] = Field{
			Name: tagValue,
		}
	}

	return fields
}

// Struct validates a struct's contents according to the rules set in the "validate" tag, and
// returns all the data needed by the template to render the form: the original submitted values,
// validation error messages, and the field names.
func Struct[T any](s T, validatorMaps ...ValidatorMap) (map[string]Field, error) {
	err := validate.Struct(s)

	// No errors and no possible additional validation functions to run means we're done here.
	if err == nil && len(validatorMaps) == 0 {
		return map[string]Field{}, nil
	}

	var errorMessages validator.ValidationErrorsTranslations
	if err != nil {
		err := err.(validator.ValidationErrors) //nolint
		errorMessages = fixErrorMessages(err.Translate(translator))
	}

	fields := map[string]Field{}
	structName := reflectutil.GetType(s).Name()
	structValue := reflectutil.GetValue(s)

	for i := 0; i < structValue.NumField(); i++ {
		submittedValue := reflectutil.ToString(structValue.Field(i).Interface())
		fieldName := structValue.Type().Field(i).Name
		tagValue := structValue.Type().Field(i).Tag.Get("form")
		errorMessage := errorMessages[structName+"."+fieldName]

		for _, validatorMap := range validatorMaps {
			if errorMessage == "" {
				if validate, ok := validatorMap[fieldName]; ok {
					if err := validate(submittedValue); err != nil {
						errorMessage = err.Error()
					}
				}
			}
			if errorMessage != "" {
				err = ErrValidationFailed
			}
		}

		fields[fieldName] = Field{
			Error: errorMessage,
			Value: submittedValue,
			Name:  tagValue,
		}
	}

	// Even after additional validators were run there are still no errors.
	if err == nil {
		return map[string]Field{}, nil
	}

	return fields, ErrValidationFailed
}

// fixErrorMessages cleans up error messages by removing the field name from the beginning as it's
// not necessary and is neater without.
func fixErrorMessages(messages validator.ValidationErrorsTranslations) validator.ValidationErrorsTranslations {
	for key, message := range messages {
		// Depending on whether an anonymous or named struct was passed in, the field name might
		// be "StructName.FieldName", or just "FieldName", so check for that.
		fieldName := key
		if strings.Contains(key, ".") {
			_, fieldName, _ = strings.Cut(key, ".")
		}

		re := regexp.MustCompile(`^` + regexp.QuoteMeta(fieldName) + `\s*`)
		message = re.ReplaceAllString(message, "")

		if message == "is a required field" {
			message = "required field"
		}

		messages[key] = message
	}

	return messages
}

// ConfirmPassword returns a validator that ensures the two passwords match.
func ConfirmPassword(password string) func(string) error {
	return func(confirmPassword string) error {
		if password != confirmPassword {
			return errors.New("passwords must match")
		}
		return nil
	}
}

// CurrentPassword returns a validator that ensures the value is valid when verified against the
// password hash.
func CurrentPassword(hash string) func(string) error {
	return func(password string) error {
		if !stringutil.VerifyHashAndPassword(hash, password) {
			return errors.New("must be your current password")
		}
		return nil
	}
}
