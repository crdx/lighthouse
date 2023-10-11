package validate

import (
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
		fields := map[string]Field{}
		structName := reflectutil.GetName(s)

		value := reflectutil.GetValue(s)

		for i := 0; i < value.NumField(); i++ {
			submittedValue := reflectutil.ToString(value.Field(i))
			fieldName := value.Type().Field(i).Name

			fields[fieldName] = Field{
				Error: errorMessages[structName+"."+fieldName],
				Value: submittedValue,
			}
		}

		return fields, true
	}
}
