package utils

// import (
// 	"errors"
// 	// "github.com/inpime/fader/utils/sdata"
// 	// "log"
// 	// "github.com/imdario/mergo"
// 	"reflect"
// )

// func resolveValues(dst, src interface{}) (vdst reflect.Value, vsrc reflect.Value, err error) {
// 	if dst == nil || src == nil {
// 		err = errors.New("nil args")
// 		return
// 	}

// 	// destination

// 	vdst = reflect.ValueOf(dst)

// 	if vdst.Kind() == reflect.Ptr {
// 		vdst = vdst.Elem()
// 	}

// 	if vdst.Kind() != reflect.Struct &&
// 		vdst.Kind() != reflect.Slice &&
// 		vdst.Kind() != reflect.Map {

// 		err = errors.New("not supported")

// 		return
// 	}

// 	// source

// 	vsrc = reflect.ValueOf(src)

// 	if vsrc.Kind() == reflect.Ptr {
// 		vsrc = vsrc.Elem()
// 	}

// 	if vdst.Type() != vsrc.Type() {
// 		err = errors.New("diff args types")
// 		return
// 	}

// 	return
// }

// func appendOrReplace(vdst, vsrc reflect.Value) error {

// 	switch vdst.Kind() {
// 	case reflect.Struct:
// 		for i := 0; i < vdst.NumField(); i++ {
// 			if err := appendOrReplace(vdst.Field(i), vsrc.Field(i)); err != nil {
// 				return err
// 			}
// 		}
// 	case reflect.Map:
// 		if !isEmptyValue(vdst) && !isEmptyValue(vsrc) {
// 			for _, key := range vsrc.MapKeys() {
// 				vdst.SetMapIndex(key, vsrc.MapIndex(key)) // if error "panic: assignment to entry in nil map".
// 				// because map not init, please make(map[...]...)
// 			}
// 		}
// 	case reflect.Slice:
// 		for i := 0; i < vsrc.Len(); i++ {
// 			vdst.Set(reflect.Append(vdst, vsrc.Index(i)))
// 		}
// 	case reflect.Ptr:
// 		if isEmptyValue(vsrc) {
// 			return nil
// 		}

// 		return appendOrReplace(vdst.Elem(), vsrc.Elem())
// 	default:
// 		if vdst.CanSet() && !isEmptyValue(vsrc) {
// 			vdst.Set(vsrc)
// 		}
// 	}

// 	return nil
// }

func AppendOrReplace(d, s interface{}) error {
	return Merge(d, s)
}

// // func AppendOrReplace(d, s interface{}) error {
// // 	vd, vs, err := resolveValues(d, s)

// // 	if err != nil {
// // 		return err
// // 	}

// // 	return appendOrReplace(vd, vs)
// // }

// // From src/pkg/encoding/json.
// func isEmptyValue(v reflect.Value) bool {
// 	switch v.Kind() {
// 	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
// 		return v.Len() == 0
// 	case reflect.Bool:
// 		return !v.Bool()
// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 		return v.Int() == 0
// 	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
// 		return v.Uint() == 0
// 	case reflect.Float32, reflect.Float64:
// 		return v.Float() == 0
// 	case reflect.Interface, reflect.Ptr:
// 		return v.IsNil()
// 	}
// 	return false
// }
