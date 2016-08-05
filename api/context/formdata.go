package context

// import "github.com/inpime/fader/utils/sdata"
import "github.com/inpime/sdata"

// BindFormToMap returns the form field values for the provided names.
func (c Context) BindFormToMap(fieldNames ...string) *sdata.StringMap {
	m := sdata.NewStringMap()

	for _, name := range fieldNames {
		m.Set(name, c.FormValue(name))
	}

	return m
}
