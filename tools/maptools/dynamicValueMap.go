package maptools

import (
	"fmt"
	"strconv"
)

type DynamicValuedMap struct {
	Params map[string]interface{} `json:"params"`
}

func NewDynamicValuedMap(m map[string]interface{}) *DynamicValuedMap {
	return &DynamicValuedMap{m}
}

func (userSubReq *DynamicValuedMap) GetFloat64Array(key string) (value []float64, err error) {
	v, ok := userSubReq.Params[key]
	if !ok {
		return nil, fmt.Errorf("value not found")
	}

	value = nil
	arr, ok := v.([]interface{})

	if !ok {
		return nil, fmt.Errorf("value to []interface{} failed")
	}

	value = make([]float64, 0)
	for _, item := range arr {
		_r, ok := getFloat(item)
		if !ok {
			return nil, fmt.Errorf("%v to float failed", item)
		}

		value = append(value, _r)
	}

	return value, nil
}

func (userSubReq *DynamicValuedMap) GetInt(key string, def ...int64) (value int64, ok bool) {
	value = 0
	ok = false

	var defVal int64 = 0
	if len(def) > 0 {
		defVal = def[0]
	}

	v, ok := userSubReq.Params[key]
	if !ok {
		return defVal, false
	}

	return getInt(v)
}

func (userSubReq *DynamicValuedMap) GetFloat(key string, def ...float64) (value float64, ok bool) {

	value = 0.0
	ok = false

	defVal := 0.0
	if len(def) > 0 {
		defVal = def[0]
	}

	v, ok := userSubReq.Params[key]
	if !ok {
		return defVal, false
	}

	return getFloat(v)
}

func (userSubReq *DynamicValuedMap) GetString(key string, def ...string) (value string, ok bool) {

	value = ""
	ok = false

	defVal := ""
	if len(def) > 0 {
		defVal = def[0]
	}

	v, ok := userSubReq.Params[key]
	if !ok {
		return defVal, false
	}

	switch v.(type) {
	case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		value = fmt.Sprintf("%d", v)
		ok = true
	case string:
		value = fmt.Sprintf("%s", v)
		ok = true
	case float64:
		value = strconv.FormatFloat(v.(float64), 'f', -1, 64)
		ok = true
	case float32:
		value = strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32)
		ok = true
	case bool:
		if v.(bool) {
			value = "true"
		} else {
			value = "false"
		}
		ok = true
	}

	return value, ok
}

func getFloat(v interface{}) (value float64, ok bool) {

	value = 0.0
	ok = false

	switch v.(type) {
	case int:
		value = float64(v.(int))
		ok = true
	case int64:
		value = float64(v.(int64))
		ok = true
	case int32:
		value = float64(v.(int32))
		ok = true
	case int16:
		value = float64(v.(int16))
		ok = true
	case int8:
		value = float64(v.(int8))
		ok = true
	case uint:
		value = float64(v.(uint))
		ok = true
	case uint64:
		value = float64(v.(uint64))
		ok = true
	case uint32:
		value = float64(v.(uint32))
		ok = true
	case uint16:
		value = float64(v.(uint16))
		ok = true
	case uint8:
		value = float64(v.(uint8))
		ok = true
	case string:
		s := v.(string)
		parsedValue, e := strconv.ParseFloat(s, 64)
		if e == nil {
			value = parsedValue
			ok = true
		}
	case float64:
		value = v.(float64)
		ok = true
	case float32:
		value = float64(v.(float32))
		ok = true
	case bool:
		if v.(bool) {
			value = 1.0
		} else {
			value = 0.0
		}
		ok = true

	}

	return value, ok
}

func getInt(v interface{}) (value int64, ok bool) {
	value = 0
	ok = false

	switch v.(type) {
	case int:
		value = int64(v.(int))
		ok = true
	case int64:
		value = v.(int64)
		ok = true
	case int32:
		value = int64(v.(int32))
		ok = true
	case int16:
		value = int64(v.(int16))
		ok = true
	case int8:
		value = int64(v.(int8))
		ok = true
	case uint:
		value = int64(v.(uint))
		ok = true
	case uint64:
		value = int64(v.(uint64))
		ok = true
	case uint32:
		value = int64(v.(uint32))
		ok = true
	case uint16:
		value = int64(v.(uint16))
		ok = true
	case uint8:
		value = int64(v.(uint8))
		ok = true
	case string:
		s := v.(string)
		parsedValue, e := strconv.ParseInt(s, 10, 64)
		if e == nil {
			value = parsedValue
			ok = true
		}
	case float64:
		value = int64(v.(float64))
		ok = true
	case float32:
		value = int64(v.(float32))
		ok = true
	case bool:
		if v.(bool) {
			value = 1
		} else {
			value = 0
		}
		ok = true

	}

	return value, ok
}
