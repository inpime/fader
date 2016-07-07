package standard

import (
	"encoding/gob"
	"fmt"
	"gopkg.in/go-playground/validator.v8"
	"utils"
)

var V = validator.New(&validator.Config{})

func init() {
	gob.Register(FIeld{})
	gob.Register(ValidatorData{})
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
				NewFIeld(fieldName,
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

func (v ValidatorData) Get(fieldName string) FIeld {
	field := v.ErrorMessages.Get(fieldName)

	if field, valid := field.(FIeld); valid {
		return field
	}

	return FIeld{
		FieldName: fieldName,
		Value:     v.Data.Get(fieldName),
	}
}

func (v ValidatorData) Messages() utils.M {
	return v.ErrorMessages
}

//

// NewFIeld информация об ошибке, только для одного поля (массив содержит только одну ошибку с пустым "" ключем)
func NewFIeld(fieldName string, value interface{}, err error) FIeld {
	if verr, ok := err.(validator.ValidationErrors); ok {
		if ferr, exists := verr[""]; exists {

			return FIeld{
				IsError:   true,
				Tag:       ferr.Tag,
				Param:     ferr.Param,
				FieldName: fieldName,
				Value:     value,
			}
		}
	}

	return FIeld{}
}

type FIeld struct {
	IsError   bool
	Tag       string
	Param     string
	FieldName string
	Value     interface{}
}

func (f FIeld) FormatMessage() string {
	return "Field validation for %q failed on the %q tag"
}

func (f FIeld) String() string {
	if !f.IsError {
		return ""
	}
	return fmt.Sprintf(f.FormatMessage(), f.FieldName, f.Tag+" "+f.Param)
}
