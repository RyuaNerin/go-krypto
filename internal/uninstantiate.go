package internal

import (
	"reflect"
)

func SetZero(v interface{}) {
	p := reflect.ValueOf(v)
	for p.Kind() == reflect.Ptr {
		p = p.Elem()
	}
	switch p.Kind() {
	case reflect.Slice:
		l := p.Len()
		for idx := 0; idx < l; idx++ {
			pidx := p.Index(idx)
			pidx.Set(reflect.Zero(pidx.Type()))
		}

	case reflect.Array:
		l := p.Len()
		for idx := 0; idx < l; idx++ {
			pidx := p.Index(idx)
			pidx.Set(reflect.Zero(pidx.Type()))
		}

	case reflect.Map:
		for _, key := range p.MapKeys() {
			p.SetMapIndex(key, reflect.Value{})
		}

	case reflect.Struct:
		p.Set(reflect.Zero(p.Type()))
	}
}
