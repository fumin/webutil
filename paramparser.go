package webutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
)

// A ParamParser is a helper to extract params with types other than strings.
type ParamParser struct {
	R   *http.Request
	W   http.ResponseWriter
	Err error
}

const errInt int64 = math.MinInt64

func (parser *ParamParser) RequiredIntParam(key string) int64 {
	if parser.Err != nil {
		return errInt
	}
	v := parser.R.FormValue(key)
	if v == "" {
		errMsg := fmt.Sprintf("missing param: %v", key)
		parser.Err = errors.New(errMsg)
		return errInt
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("Wrong integer format: %v", v)
		parser.Err = errors.New(errMsg)
		return errInt
	}
	return i
}

func (parser *ParamParser) OptionalIntParam(key string, defaultVal int64) int64 {
	if parser.Err != nil {
		return errInt
	}
	v := parser.R.FormValue(key)
	if v == "" {
		return defaultVal
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("Wrong integer format: %v", v)
		parser.Err = errors.New(errMsg)
		return errInt
	}
	return i
}

const errFloat float64 = math.SmallestNonzeroFloat64

func (parser *ParamParser) RequiredFloatParam(key string) float64 {
	if parser.Err != nil {
		return errFloat
	}
	v := parser.R.FormValue(key)
	if v == "" {
		errMsg := fmt.Sprintf("missing param: %v", key)
		parser.Err = errors.New(errMsg)
		return errFloat
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		errMsg := fmt.Sprintf("Wrong float format: %v", v)
		parser.Err = errors.New(errMsg)
		return errFloat
	}
	return f
}

func (parser *ParamParser) OptionalFloatParam(key string, defaultVal float64) float64 {
	if parser.Err != nil {
		return errFloat
	}
	v := parser.R.FormValue(key)
	if v == "" {
		return defaultVal
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		errMsg := fmt.Sprintf("Wrong float format: %v", v)
		parser.Err = errors.New(errMsg)
		return errFloat
	}
	return f
}

func JsonResp(w http.ResponseWriter, o interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(o)
}
