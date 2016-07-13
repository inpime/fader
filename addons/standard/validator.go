package standard

import (
	"encoding/gob"
	"fmt"
	"github.com/inpime/fader/utils/sdata"
	"gopkg.in/go-playground/validator.v8"
)

var V = validator.New(&validator.Config{})

func init() {
	gob.Register(FIeld{})
	gob.Register(ValidatorData{})
}

func NewValidatorData() *ValidatorData {
	return &ValidatorData{
		Data:          sdata.NewStringMap(),
		Rules:         sdata.NewStringMap(),
		ErrorMessages: sdata.NewStringMap(),
	}
}

type ValidatorData struct {
	Data  *sdata.StringMap
	Rules *sdata.StringMap

	ErrorMessages *sdata.StringMap
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
	v.ErrorMessages = sdata.NewStringMap() // clear old error messages

	for _, fieldName := range v.Rules.Keys() {
		rule := v.Rules.String(fieldName)
		if err := V.Field(v.Data.GetOrNil(fieldName), rule); err != nil {

			v.ErrorMessages.Set(fieldName,
				NewFIeld(fieldName,
					v.Data.GetOrNil(fieldName),
					err,
				))
		}
	}

	if v.ErrorMessages.Size() > 0 {
		return false
	}

	return true
}

func (v ValidatorData) Get(fieldName string) FIeld {
	field := v.ErrorMessages.GetOrNil(fieldName)

	if field, valid := field.(FIeld); valid {
		return field
	}

	return FIeld{
		FieldName: fieldName,
		Value:     v.Data.GetOrNil(fieldName),
	}
}

func (v ValidatorData) Messages() *sdata.StringMap {
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
