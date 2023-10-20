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

// Struct validates a struct's contents according to the rules set in the "validate" tag, and
// returns all the data needed by the template to render the form: the original submitted values
// and any validation error messages.
func Struct[T any](s T) (map[string]Field, bool) {
	if err := validate.Struct(s); err == nil {
		return nil, false
	} else {
		err := err.(validator.ValidationErrors) //nolint
		errorMessages := err.Translate(translator)
		errorMessages = unProperCaseMessages(errorMessages)

		fields := map[string]Field{}
		structName := reflectutil.GetType(s).Name()
		structValue := reflectutil.GetValue(s)

		for i := 0; i < structValue.NumField(); i++ {
			submittedValue := reflectutil.ToString(structValue.Field(i).Interface())
			fieldName := structValue.Type().Field(i).Name

			fields[fieldName] = Field{
				Error: errorMessages[structName+"."+fieldName],
				Value: submittedValue,
			}
		}

		return fields, true
	}
}

// unProperCaseMessages takes a map of error messages and converts ProperCased field names in the
// message to more readable ones. For example, "GracePeriod" turns into "Grace Period".
func unProperCaseMessages(messages validator.ValidationErrorsTranslations) validator.ValidationErrorsTranslations {
	for key, message := range messages {
		// Depending on whether an anonymous or named struct was passed in, the field name might
		// be "StructName.FieldName", or just "FieldName", so check for that.
		fieldName := key
		if strings.Contains(key, ".") {
			_, fieldName, _ = strings.Cut(key, ".")
		}

		// Find all propercased words in the field name e.g. GracePeriod is "Grace" and "Period".
		words := regexp.MustCompile(`([A-Z][a-z]*)`).FindAllString(fieldName, -1)

		if words != nil {
			// Prepare a regex of the field name where it appears as a whole word, to prevent any
			// erroneous replacements if the field name is a very simple and common sequence.
			re := regexp.MustCompile(`\b` + regexp.QuoteMeta(fieldName) + `\b`)

			messages[key] = re.ReplaceAllString(message, strings.Join(words, " "))
		}
	}

	return messages
}
