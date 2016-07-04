package standard

import (
	"encoding/gob"
	"fmt"
	"gopkg.in/go-playground/validator.v8"
	"utils"
)

var V = validator.New(&validator.Config{})

func init() {
	gob.Register(FieldError{})
}

func NewValidatorData() *ValidatorData {
	return &ValidatorData{
		Data:          utils.M{},
		Rules:         utils.M{},
		ErrorMessages: utils.M{},
	}
}

type ValidatorData struct {
	Data  utils.M
	Rules utils.M

	ErrorMessages utils.M
}

func (v *ValidatorData) SetData(data interface{}) *ValidatorData {
	v.Data.LoadFrom(data)
	return v
}

func (v *ValidatorData) AddRule(fieldName, ruleStr string) *ValidatorData {
	v.Rules.Set(fieldName, ruleStr)
	return v
}

func (v *ValidatorData) Valid() bool {
	v.ErrorMessages = utils.M{} // clear old error messages

	for fieldName, rule := range v.Rules {
		if err := V.Field(v.Data.Get(fieldName), rule.(string)); err != nil {

			v.ErrorMessages.Set(fieldName,
				NewFieldError(fieldName,
					v.Data.Get(fieldName),
					err,
				))
		}
	}

	if len(v.ErrorMessages) > 0 {
		return false
	}

	return true
}

func (v ValidatorData) Messages() utils.M {
	return v.ErrorMessages
}

//

// NewFieldError информация об ошибке, только для одного поля (массив содержит только одну ошибку с пустым "" ключем)
func NewFieldError(fieldName string, value interface{}, err error) FieldError {
	if verr, ok := err.(validator.ValidationErrors); ok {
		if ferr, exists := verr[""]; exists {

			return FieldError{
				Tag:       ferr.Tag,
				Param:     ferr.Param,
				FieldName: fieldName,
				Value:     value,
			}
		}
	}

	return FieldError{}
}

type FieldError struct {
	Tag       string
	Param     string
	FieldName string
	Value     interface{}
}

func (f FieldError) FormatMessage() string {
	return "Field validation for %q failed on the %q tag"
}

func (f FieldError) String() string {
	return fmt.Sprintf(f.FormatMessage(), f.FieldName, f.Tag)
}
