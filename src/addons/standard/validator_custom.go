package standard

import (
	"gopkg.in/go-playground/validator.v8"
	"reflect"
	"strconv"
)

func init() {
	V.RegisterValidation("atoigt", AtoiGT)
}

func AtoiGT(v *validator.Validate, topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldtype reflect.Type, fieldKind reflect.Kind, param string) bool {

	if fieldKind != reflect.String {
		return false
	}

	i, err := strconv.Atoi(field.String())
	if err != nil {
		return false
	}

	vo := reflect.ValueOf(i)

	return validator.IsGt(v, topStruct, field, vo, vo.Type(), vo.Kind(), param)
}
