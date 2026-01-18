package binder

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

const (
	Path   = "path"
	Query  = "query"
	Header = "header"
)

type Provider func(r *http.Request, key string) ([]string, bool)

type Binder struct {
	providers map[string]Provider
}

func NewDefault() *Binder {
	b := &Binder{
		providers: make(map[string]Provider),
	}

	b.Register(Path, func(r *http.Request, key string) ([]string, bool) {
		val := r.PathValue(key)
		if val == "" {
			return nil, false
		}
		return []string{val}, true
	})

	b.Register(Query, func(r *http.Request, key string) ([]string, bool) {
		if !r.URL.Query().Has(key) {
			return nil, false
		}
		return r.URL.Query()[key], true
	})

	b.Register(Header, func(r *http.Request, key string) ([]string, bool) {
		if vals, ok := r.Header[http.CanonicalHeaderKey(key)]; ok && len(vals) > 0 {
			return vals, true
		}
		return nil, false
	})

	return b
}

func (b *Binder) Register(tagName string, p Provider) {
	b.providers[tagName] = p
}

func (b *Binder) Bind(r *http.Request, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("binder: bind target must be a non-nil pointer")
	}
	return b.bindRecursive(r, v.Elem())
}

func (b *Binder) bindRecursive(r *http.Request, v reflect.Value) error {
	t := v.Type()

	if v.Kind() == reflect.Struct {
		for i := 0; i < v.NumField(); i++ {
			fieldVal := v.Field(i)
			fieldType := t.Field(i)

			if !fieldVal.CanSet() {
				continue
			}

			if fieldVal.Kind() == reflect.Struct {
				if err := b.bindRecursive(r, fieldVal); err != nil {
					return err
				}
				continue
			}

			if fieldVal.Kind() == reflect.Pointer && fieldVal.Type().Elem().Kind() == reflect.Struct {
				if fieldVal.IsNil() {
					fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
				}
				if err := b.bindRecursive(r, fieldVal.Elem()); err != nil {
					return err
				}
				continue
			}

			if err := b.bindField(r, fieldVal, fieldType); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *Binder) bindField(r *http.Request, fieldVal reflect.Value, fieldType reflect.StructField) error {
	tagStr := string(fieldType.Tag)

	values, found := b.processTagsLazy(r, tagStr)

	if found && len(values) > 0 {
		if fieldVal.Kind() == reflect.Slice {
			if fieldVal.Type().Elem().Kind() == reflect.Uint8 {
				return nil
			}
			return setSlice(fieldVal, values)
		}

		return setScalar(fieldVal, values[0])
	}

	return nil
}

func (b *Binder) processTagsLazy(r *http.Request, tagStr string) ([]string, bool) {
	for tagStr != "" {
		i := 0
		for i < len(tagStr) && tagStr[i] == ' ' {
			i++
		}
		tagStr = tagStr[i:]
		if tagStr == "" {
			break
		}

		j := 0
		for j < len(tagStr) && tagStr[j] > ' ' && tagStr[j] != ':' && tagStr[j] != '"' && tagStr[j] != 0x7f {
			j++
		}
		if j == 0 || j+1 >= len(tagStr) || tagStr[j] != ':' || tagStr[j+1] != '"' {
			break
		}
		key := tagStr[:j]
		tagStr = tagStr[j+1:]

		j = 1
		for j < len(tagStr) && tagStr[j] != '"' {
			if tagStr[j] == '\\' {
				j++
			}
			j++
		}
		if j >= len(tagStr) {
			break
		}
		qvalue := tagStr[:j+1]
		tagStr = tagStr[j+1:]

		tagName, err := strconv.Unquote(qvalue)
		if err != nil {
			continue
		}

		if provider, exists := b.providers[key]; exists {
			if vals, ok := provider(r, tagName); ok {
				return vals, true
			}
		}
	}
	return nil, false
}

func setScalar(v reflect.Value, s string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)
	}
	return nil
}

func setSlice(v reflect.Value, values []string) error {
	sliceType := v.Type()

	newSlice := reflect.MakeSlice(sliceType, len(values), len(values))

	for i, s := range values {
		if err := setScalar(newSlice.Index(i), s); err != nil {
			return fmt.Errorf("slice index %d: %w", i, err)
		}
	}

	v.Set(newSlice)
	return nil
}
