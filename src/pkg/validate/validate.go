package validate

import (
	"regexp"
	"strings"

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
	if validate == nil {
		panic("unable to build validator")
	}

	lo.Must0(enTranslations.RegisterDefaultTranslations(validate, translator))
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
		errorMessages = removeFieldName(errorMessages)

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

// removeFieldName removes the field name from the beginning of messages as it's not necessary and
// is neater without.
func removeFieldName(messages validator.ValidationErrorsTranslations) validator.ValidationErrorsTranslations {
	for key, message := range messages {
		// Depending on whether an anonymous or named struct was passed in, the field name might
		// be "StructName.FieldName", or just "FieldName", so check for that.
		fieldName := key
		if strings.Contains(key, ".") {
			_, fieldName, _ = strings.Cut(key, ".")
		}

		re := regexp.MustCompile(`^` + regexp.QuoteMeta(fieldName) + `\s*`)
		messages[key] = re.ReplaceAllString(message, "")
	}

	return messages
}
