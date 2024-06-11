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

func SetZeroResurvie(v interface{}) {
	hit := make(map[uintptr]bool)
	uninstantiateField(v, hit)
	SetZero(v)
}

func uninstantiateField(value interface{}, visited map[uintptr]bool) {
	v := reflect.ValueOf(value)
	{
		ptr := reflect.ValueOf(value).Pointer()
		if _, ok := visited[ptr]; ok {
			return
		}
		visited[ptr] = true
	}

	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if !v.IsValid() {
		return
	}
	vt := v.Type()

	for fieldIdx := 0; fieldIdx < v.NumField(); fieldIdx++ {
		vtField := vt.Field(fieldIdx)
		if !vtField.IsExported() {
			continue
		}

		field := v.Field(fieldIdx)
		switch field.Kind() {
		case reflect.Pointer:
			uninstantiateField(field.Interface(), visited)
		case reflect.Slice:
			l := field.Len()
			for idx := 0; idx < l; idx++ {
				field.Index(idx).Set(reflect.Zero(field.Index(idx).Type()))
			}
		case reflect.Map:
			for _, key := range field.MapKeys() {
				field.SetMapIndex(key, reflect.Value{})
			}
		}
	}
}
