package context

import "utils"

// BindFormToMap helper function for fast bind forms
func (c Context) BindFormToMap(fieldNames ...string) utils.M {
	m := utils.Map()

	for _, name := range fieldNames {
		m.Set(name, c.FormValue(name))
	}

	return m
}
