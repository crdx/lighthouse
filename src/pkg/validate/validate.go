package validate

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"crdx.org/lighthouse/util/reflectutil"
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

var validate *validator.Validate
var translator universalTranslator.Translator

func init() {
	translator, _ = universalTranslator.New(enLocale.New()).GetTranslator("en")
	validate = validator.New(validator.WithRequiredStructEnabled())

	lo.Must0(enTranslations.RegisterDefaultTranslations(validate, translator))

	Register("timezone", "must be valid", func(value string) bool {
		_, err := time.LoadLocation(value)
		return err == nil
	})

	Register("mailaddr", `must be in the format "xxx <yyy>"`, func(value string) bool {
		return regexp.MustCompile("^.* <.*>$").Match([]byte(value))
	})

	Register("scan_interval", `must be between 1 and 30 minutes`, func(value string) bool {
		n, err := strconv.Atoi(value)
		if err != nil {
			return false
		}
		return n >= 1 && n <= 30
	})

	Register("grace_period", `must be between 1 and 60 minutes`, func(value string) bool {
		n, err := strconv.Atoi(value)
		if err != nil {
			return false
		}
		return n >= 1 && n <= 60
	})
}

// Register registers a custom validation function.
func Register(name string, message string, f func(string) bool) {
	lo.Must0(validate.RegisterValidation(name, func(field validator.FieldLevel) bool {
		return f(field.Field().String())
	}, false))

	registerFn := func(translator universalTranslator.Translator) error {
		return translator.Add(name, message, true)
	}

	translateFn := func(translator universalTranslator.Translator, _ validator.FieldError) string {
		return lo.Must(translator.T(name))
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
func Struct[T any](s T) (map[string]Field, bool) {
	if err := validate.Struct(s); err == nil {
		return nil, false
	} else {
		err := err.(validator.ValidationErrors) //nolint
		errorMessages := err.Translate(translator)
		errorMessages = fixErrorMessages(errorMessages)

		fields := map[string]Field{}
		structName := reflectutil.GetType(s).Name()
		structValue := reflectutil.GetValue(s)

		for i := 0; i < structValue.NumField(); i++ {
			submittedValue := reflectutil.ToString(structValue.Field(i).Interface())
			fieldName := structValue.Type().Field(i).Name
			tagValue := structValue.Type().Field(i).Tag.Get("form")

			fields[fieldName] = Field{
				Error: errorMessages[structName+"."+fieldName],
				Value: submittedValue,
				Name:  tagValue,
			}
		}

		return fields, true
	}
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
