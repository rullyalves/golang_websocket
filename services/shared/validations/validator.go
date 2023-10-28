package validations

import (
	"context"
	"errors"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	valid "github.com/go-playground/validator/v10"
	ptTranslations "github.com/go-playground/validator/v10/translations/pt_BR"
)

var (
	uni       *ut.UniversalTranslator
	validator *valid.Validate
)

type ErrorItem struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func InitTranslations() {

	validator = valid.New()

	ptBrLang := pt_BR.New()

	uni = ut.New(ptBrLang, ptBrLang)

	ptBrTranslator, _ := uni.GetTranslator("pt_BR")

	_ = ptTranslations.RegisterDefaultTranslations(validator, ptBrTranslator)

}

func ValidateErrors(ctx context.Context, language string, value any) []ErrorItem {

	validationErrors := validate(ctx, value)

	if len(validationErrors) == 0 {
		return []ErrorItem{}
	}

	var errorList []ErrorItem

	translator, _ := uni.GetTranslator(language)

	for _, item := range validationErrors {
		errorList = append(errorList, ErrorItem{Field: item.Field(), Message: item.Translate(translator)})
	}

	return errorList
}

func validate(ctx context.Context, value any) valid.ValidationErrors {

	err := validator.StructCtx(ctx, value)

	var validationErr valid.ValidationErrors

	if errors.As(err, &validationErr) {
		return validationErr
	}

	return valid.ValidationErrors{}
}
