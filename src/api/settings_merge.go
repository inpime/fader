package api

import (
	"errors"
	"fmt"
	"reflect"
)

func resolve(dst, src interface{}) (
	vdst reflect.Value,
	vsrc reflect.Value,
	tdst reflect.Type,
	tsrc reflect.Type,
	err error,
) {
	tdst = reflect.TypeOf(dst)
	tsrc = reflect.TypeOf(src)
	vdst = reflect.ValueOf(dst)
	vsrc = reflect.ValueOf(src)

	vsrc = reflect.Indirect(vsrc)
	vdst = reflect.Indirect(vdst)

	if vdst.Type() != vsrc.Type() {
		err = errors.New("diff args types")
		return
	}

	switch vdst.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Map:

	default:
		err = errors.New("not supported")

		return
	}

	return
}

func Merge(d, s interface{}) error {
	// vdst, vsrc, tdst, tsrc, err := resolve(d, s)
	vdst, vsrc, _, _, err := resolve(d, s)

	if err != nil {
		return err
	}

	return merge(vdst, vsrc, 1 /*, tdst, tsrc*/)
}

func merge(vdst reflect.Value,
	vsrc reflect.Value,
	dep int,
	/*tdst reflect.Type,
	tsrc reflect.Type*/) error {

	// log.Printf(strings.Repeat("  ", dep+1)+"%v %v (can set %v, can addr %v)", vdst.Kind(), vsrc.Kind(), vdst.CanSet(), vdst.CanAddr())

	switch vdst.Kind() {
	case reflect.Struct:
		for i := 0; i < vdst.NumField(); i++ {
			_vdst := vdst.Field(i)
			_vsrc := vsrc.Field(i)
			// _tdst := _vdst.Type()
			// _tsrc := _vsrc.Type()

			if err := merge(_vdst, _vsrc, dep+1 /*, _tdst, _tsrc*/); err != nil {
				return err
			}
		}
	case reflect.Slice:
		for i := 0; i < vsrc.Len(); i++ {
			vdst.Set(reflect.Append(vdst, vsrc.Index(i)))
		}
	case reflect.Map:
		for _, key := range vsrc.MapKeys() {
			_vsrc := vsrc.MapIndex(key)

			if isSimpleType(_vsrc) {
				vdst.SetMapIndex(key, _vsrc)
			} else {

				switch _vsrc.Kind() {
				case reflect.Interface:
					_vsrc = reflect.ValueOf(_vsrc.Interface())

					switch _vsrc.Kind() {
					case reflect.String,
						reflect.Int,
						reflect.Int16,
						reflect.Int32,
						reflect.Int64,
						reflect.Int8,
						reflect.Float32,
						reflect.Float64,
						reflect.Bool,
						reflect.Chan,
						reflect.Func,
						reflect.Uint,
						reflect.Uint8,
						reflect.Uint16,
						reflect.Uint32,
						reflect.Uint64:

						vdst.SetMapIndex(key, _vsrc)
					case reflect.Slice:
						_vdst := vdst.MapIndex(key)

						if _vdst.IsValid() {
							// log.Printf(strings.Repeat("  ", dep+1)+"%v: %v %v (can set %v) %#v %#v", key, _vdst.Type(), _vsrc.Type(), _vdst.CanSet(), vdst.MapIndex(key).Interface(), _vsrc.Interface())
						} else {
							// log.Printf(strings.Repeat("  ", dep+1)+"%v: %v %v (can set %v) %#v %#v", key, "nil", _vsrc.Type(), _vdst.CanSet(), "nil", _vsrc.Interface())
						}

						if !_vdst.CanSet() {
							_vdst = reflect.New(_vsrc.Type()).Elem()
							_len := _vsrc.Len()

							if vdst.MapIndex(key).IsValid() {

								_len += reflect.ValueOf(vdst.MapIndex(key).Interface()).Len()

								_vdst.Set(reflect.ValueOf(vdst.MapIndex(key).Interface())) // interface to _vsrc.Type()
							}
						}

						for i := 0; i < _vsrc.Len(); i++ {
							_vdst.Set(reflect.Append(_vdst, _vsrc.Index(i)))
						}

						vdst.SetMapIndex(key, reflect.Indirect(_vdst))
					case reflect.Map:
						_vdst := vdst.MapIndex(key)

						if !_vdst.CanSet() {
							_vdst = reflect.New(_vsrc.Type()).Elem()
							_vdst.Set(reflect.MakeMap(_vsrc.Type()))

							// _vdst.CanSet()
							if vdst.MapIndex(key).IsValid() {
								// log.Printf(strings.Repeat("  ", dep+1)+"%v: %v %v (can set %v) %#v", key, _vdst.Type(), _vsrc.Type(), _vdst.CanSet(), vdst.MapIndex(key).Interface(), _vsrc.Interface())
								_vdst.Set(reflect.ValueOf(vdst.MapIndex(key).Interface())) // interface to _vsrc.Type()
							}

						}

						if err := merge(_vdst, _vsrc, dep+1 /*, _tdst, _tsrc*/); err != nil {
							return err
						}

						vdst.SetMapIndex(key, reflect.Indirect(_vdst))
					default:
						if _vsrc.CanInterface() {
							return fmt.Errorf("new supported type %s (%T)", _vsrc.Kind().String(), _vsrc.Interface())
						}
						return fmt.Errorf("new supported type %s (can interface false)", _vsrc.Kind().String())
					}
				case reflect.Map, reflect.Struct:
					_vdst := vdst.MapIndex(key)

					if !_vdst.CanSet() {
						_vdst = reflect.New(_vsrc.Type()).Elem()

						if vdst.MapIndex(key).IsValid() {
							_vdst.Set(vdst.MapIndex(key))
						}
					}

					if err := merge(_vdst, _vsrc, dep+1 /*, _tdst, _tsrc*/); err != nil {
						return err
					}

					vdst.SetMapIndex(key, reflect.Indirect(_vdst))
				}
			}
		}
	default:
		if vdst.CanSet() {
			vdst.Set(vsrc)
		}
	}

	return nil
}

func equalValues(v1, v2 reflect.Value) bool {
	if v1.Type() == v2.Type() {
		return true
	}

	return false
}

func isSimpleType(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Map, reflect.Struct, reflect.Slice, reflect.Interface:
		return false
	default:
		return true
	}
}
