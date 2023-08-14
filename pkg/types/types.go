package types

import "errors"

var (
	ErrInvalidBoolValue   = errors.New("invalid value in boolean field")
	ErrInvalidNumberValue = errors.New("invalid value in number field")
	ErrInvalidNullValue   = errors.New("invalid value in null field")
	ErrInvalidStringValue = errors.New("invalid value in string field")
	ErrEmptyList          = errors.New("list is empty")
	ErrInvalidAttribute   = errors.New("unknown attribute type")
)

type Attribute struct {
	BOOL *string               `type:"boolean"`
	NULL *string               `type:"boolean"`
	S    *string               `type:"string"`
	N    *string               `type:"string"`
	L    []*Attribute          `type:"list"`
	M    map[string]*Attribute `type:"map"`
}

type AttributeMap = map[string]*Attribute
type JSON = map[string]interface{}
