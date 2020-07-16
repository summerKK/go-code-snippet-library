package app

import (
	"strings"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidatorError struct {
	Key     string
	Message string
}

func (v *ValidatorError) Error() string {
	return v.Message
}

type ValidatorErrors []*ValidatorError

func (v ValidatorErrors) Error() string {
	return strings.Join(v.Errors(), ",")
}

func (v ValidatorErrors) Errors() []string {
	var errors []string
	for _, err := range v {
		errors = append(errors, err.Error())
	}

	return errors
}

func BindAndValid(c *gin.Context, v interface{}) (bool, ValidatorErrors) {
	var errs ValidatorErrors
	err := c.ShouldBind(v)
	if err != nil {
		v := c.Value("trans")
		trans, _ := v.(ut.Translator)
		validErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return true, nil
		}

		for k, v := range validErrors.Translate(trans) {
			errs = append(errs, &ValidatorError{
				Key:     k,
				Message: v,
			})
		}

		return true, errs
	}

	return false, nil
}
