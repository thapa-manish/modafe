package encoder

import (
	"encoding/json"
	"fmt"
	"modafe/pkg/types"
	"reflect"
	"strings"
)

type Encoder struct{}

type EncoderOption = func(*Encoder)

func NewEncoder(opts ...EncoderOption) *Encoder {
	en := &Encoder{}
	for _, op := range opts {
		op(en)
	}
	return en
}

func (en *Encoder) EncodeString(input string) (types.AttributeMap, error) {
	var inputJson types.JSON
	input = strings.ReplaceAll(input, " \"", "\"")
	input = strings.ReplaceAll(input, "\" ", "\"")
	target := make(types.AttributeMap)
	err := json.Unmarshal([]byte(input), &inputJson)
	if err != nil {
		return nil, err
	}

	for key, value := range inputJson {
		attribute, err := en.Encode(value)
		if err != nil {
			continue
		}
		target[key] = attribute
	}
	return target, nil
}

func (en *Encoder) Encode(in interface{}) (*types.Attribute, error) {
	av := &types.Attribute{}
	if err := en.encode(av, reflect.ValueOf(in)); err != nil {
		return nil, err
	}
	return av, nil
}

func (en *Encoder) encode(av *types.Attribute, v reflect.Value) error {
	v = valueElem(v)
	data, err := json.Marshal(v.Interface())
	if err == nil {
		err := json.Unmarshal(data, av)
		if err == nil {
			return nil
		}
	}

	switch v.Kind() {
	case reflect.Map:
		return en.encodeMap(av, v)
	case reflect.Slice, reflect.Array:
		return en.encodeSlice(av, v)
	default:
		return types.ErrInvalidAttribute
	}
}

func (e *Encoder) encodeMap(av *types.Attribute, v reflect.Value) error {
	tempMap := map[string]*types.Attribute{}
	for _, key := range v.MapKeys() {
		keyName := fmt.Sprint(key.Interface())
		elemVal := v.MapIndex(key)
		elem := &types.Attribute{}
		switch keyName {
		case "N", "S", "NULL", "BOOL":
			return e.encodeString(av, valueElem(elemVal), keyName)
		default:
			err := e.encode(elem, elemVal)
			if err != nil {
				continue
			}
		}

		tempMap[keyName] = elem
	}
	av.M = tempMap
	return nil
}

func (e *Encoder) encodeSlice(av *types.Attribute, v reflect.Value) error {
	if v.Kind() != reflect.Array || (v.Kind() == reflect.Array && v.Len() == 0) {
		return types.ErrEmptyList
	}

	count := 0
	for i := 0; i < v.Len(); i++ {
		elem := types.Attribute{}
		err := e.encode(&elem, v.Index(i))
		if err != nil {
			continue
		}
		av.L = append(av.L, &elem)
		count++
	}

	if count == 0 {
		return types.ErrEmptyList
	}
	return nil
}

func (e *Encoder) encodeString(av *types.Attribute, v reflect.Value, kind string) error {
	if v.Kind() != reflect.String || v.String() == "" {
		return types.ErrInvalidStringValue
	}
	s := v.String()
	switch kind {
	case "N":
		av.N = &s
	case "S":
		av.S = &s
	case "BOOL":
		av.BOOL = &s
	case "NULL":
		av.NULL = &s
	default:
		return types.ErrInvalidAttribute
	}
	return nil
}

func valueElem(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}
	return v
}
