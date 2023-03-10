package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// DefaultJSONMarshal produces pretty JSON with 2-space indentation
func DefaultJSONMarshal(v interface{}) ([]byte, error) {
	bs, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, err
	}
	return bs, nil
}

// SetIfNotDefault sets dest to the value of src if src is not the default
// value of the type.
// dest must be a pointer.
func SetIfNotDefault(src interface{}, dest interface{}) {
	switch src.(type) {
	case time.Duration:
		t := src.(time.Duration)
		if t != 0 {
			*dest.(*time.Duration) = t
		}
	case string:
		str := src.(string)
		if str != "" {
			*dest.(*string) = str
		}
	case uint64:
		n := src.(uint64)
		if n != 0 {
			*dest.(*uint64) = n
		}
	case int:
		n := src.(int)
		if n != 0 {
			*dest.(*int) = n
		}
	case float64:
		n := src.(float64)
		if n != 0 {
			*dest.(*float64) = n
		}
	case bool:
		b := src.(bool)
		if b {
			*dest.(*bool) = b
		}
	}
}

type hiddenField struct{}

func (hf hiddenField) MarshalJSON() ([]byte, error) {
	return []byte(`"XXX_hidden_XXX"`), nil
}
func (hf hiddenField) UnmarshalJSON(b []byte) error { return nil }

// DisplayJSON takes pointer to a JSON-friendly configuration struct and
// returns the JSON-encoded representation of it filtering out any struct
// fields marked with the tag `hidden:"true"`, but keeping fields marked
// with `"json:omitempty"`.
func DisplayJSON(cfg interface{}) ([]byte, error) {
	cfg = reflect.Indirect(reflect.ValueOf(cfg)).Interface()
	origStructT := reflect.TypeOf(cfg)
	if origStructT.Kind() != reflect.Struct {
		panic("the given argument should be a struct")
	}

	hiddenFieldT := reflect.TypeOf(hiddenField{})

	// create a new struct type with same fields
	// but setting hidden fields as hidden.
	finalStructFields := []reflect.StructField{}
	for i := 0; i < origStructT.NumField(); i++ {
		f := origStructT.Field(i)
		hidden := f.Tag.Get("hidden") == "true"
		if f.PkgPath != "" { // skip unexported
			continue
		}
		if hidden {
			f.Type = hiddenFieldT
		}

		// remove omitempty from tag, ignore other tags except json
		var jsonTags []string
		for _, s := range strings.Split(f.Tag.Get("json"), ",") {
			if s != "omitempty" {
				jsonTags = append(jsonTags, s)
			}
		}
		f.Tag = reflect.StructTag(fmt.Sprintf("json:\"%s\"", strings.Join(jsonTags, ",")))

		finalStructFields = append(finalStructFields, f)
	}

	// Parse the original JSON into the new
	// struct and re-convert it to JSON.
	finalStructT := reflect.StructOf(finalStructFields)
	finalValue := reflect.New(finalStructT)
	data := finalValue.Interface()
	origJSON, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(origJSON, data)
	if err != nil {
		return nil, err
	}
	return DefaultJSONMarshal(data)
}
