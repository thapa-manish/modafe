package decoder

import (
	"errors"
	"modafe/pkg/types"
	"strconv"
	"strings"
	"time"
)

type DecoderOption = func(*Decoder)
type Decoder struct{}

func NewDecoder(opts ...DecoderOption) *Decoder {
	dc := &Decoder{}
	for _, op := range opts {
		op(dc)
	}
	return dc
}

func (d *Decoder) Decode(input types.AttributeMap, output types.JSON) {
	for key, av := range input {
		if key == "" {
			continue
		}
		value, err := d.decode(av)
		if err != nil {
			continue
		}
		output[key] = value
	}
}

func (d *Decoder) decode(av *types.Attribute) (interface{}, error) {
	if av == nil {
		return nil, errors.New("invalid attribute value")
	}
	switch {
	case av.BOOL != nil:
		return d.decodeBool(strings.Trim(*av.BOOL, " "))

	case av.NULL != nil:
		value, _ := d.decodeBool(strings.Trim(*av.NULL, " "))
		if !value {
			return nil, types.ErrInvalidBoolValue
		}
		return nil, nil

	case av.S != nil:
		return d.decodeString(strings.Trim(*av.S, " "))

	case av.N != nil:
		return d.decodeNumber(strings.Trim(*av.N, " "))
	case len(av.L) > 0:
		return d.decodeList(av.L)
	case av.M != nil:
		return d.decodeMap(av.M)
	default:
		return nil, types.ErrInvalidAttribute
	}
}

func (d *Decoder) decodeBool(input string) (bool, error) {
	input = strings.ToLower(input)
	if input == "f" || input == "false" || input == "0" {
		return false, nil
	} else if input == "t" || input == "true" || input == "1" {
		return true, nil
	}
	return false, types.ErrInvalidBoolValue
}

func (d *Decoder) decodeString(input string) (interface{}, error) {
	if input == "" {
		return nil, types.ErrInvalidStringValue
	}
	parsedTIme, err := time.Parse(time.RFC3339, input)
	if err == nil {
		return parsedTIme.Unix(), nil
	}
	return input, nil
}

func (d *Decoder) decodeNumber(input string) (interface{}, error) {
	if input == "" {
		return nil, types.ErrInvalidNumberValue
	}
	intValue, err := strconv.ParseInt(input, 10, 64)
	if err == nil {
		return intValue, nil
	}
	floatValue, err := strconv.ParseFloat(input, 64)
	if err != nil {
		err = types.ErrInvalidNumberValue
	}
	return floatValue, err
}

func (d *Decoder) decodeList(input []*types.Attribute) ([]interface{}, error) {
	output := make([]interface{}, 0)
	for _, av := range input {
		value, err := d.decode(av)
		if err != nil {
			continue
		}
		output = append(output, value)
	}
	if len(output) == 0 {
		return nil, types.ErrEmptyList
	}
	return output, nil
}

func (d *Decoder) decodeMap(input types.AttributeMap) (types.JSON, error) {
	output := make(types.JSON)
	validKeys := 0
	for key, av := range input {
		if key == "" {
			continue
		}
		value, err := d.decode(av)
		if err != nil {
			continue
		}
		output[key] = value
		validKeys++
	}
	if validKeys > 0 {
		return output, nil
	}
	return nil, errors.New("empty map")
}
