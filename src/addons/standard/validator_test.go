package standard

import (
	"testing"
)

func Test_SimpleValidatorTest(t *testing.T) {
	if err := V.Field("2", "required,atoigt=1"); err != nil {
		t.Error(err)
	}

	if err := V.Field(2, "required,gt=1"); err != nil {
		t.Error(err)
	}
}
