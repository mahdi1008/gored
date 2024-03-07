package respser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const CRLF = "\r\n"

var (
	ErrInvalidType       = errors.New("invalid type")
	ErrInvalidInputData  = errors.New("invalid input data")
	ErrDataMismatch      = errors.New("data mismatch")
	ErrInvalidInputParts = errors.New("invalid input parts")
	ErrFailedExtraction  = errors.New("failed extraction")
)

type RespSerError struct {
	Func string // the failing function (RespEncode, RespDecode)
	In   string // the input
	Err  error  // the reason the conversion failed (e.g. ErrInvalidType, ErrInvalidInputParts, etc.)
}

func (e *RespSerError) Error() string {
	return "respser." + e.Func + ": " + "parsing " + e.In + ": " + e.Err.Error()
}

func (e *RespSerError) Unwrap() error { return e.Err }

func invalidTypeError(fn string, in string) *RespSerError {
	return &RespSerError{fn, in, ErrInvalidType}
}

func invalidInputDataError(fn string, in string) *RespSerError {
	return &RespSerError{fn, in, ErrInvalidInputData}
}

func dataMismatchError(fn string, in string) *RespSerError {
	return &RespSerError{fn, in, ErrDataMismatch}
}

func invalidInputPartsError(fn string, in string) *RespSerError {
	return &RespSerError{fn, in, ErrInvalidInputParts}
}

func failedExtractionError(fn string, in string) *RespSerError {
	return &RespSerError{fn, in, ErrFailedExtraction}
}

type RespEncoder interface {
	RespEncode() string
}

type SimpleString struct {
	S string
}

func (ss *SimpleString) RespEncode() string {
	return fmt.Sprintf("+%s\r\n", ss.S)
}

type ErrorString struct {
	E string
}

func (es *ErrorString) RespEncode() string {
	return fmt.Sprintf("-%s\r\n", es.E)
}

type Integer struct {
	N int
}

func (i *Integer) RespEncode() string {
	nStr := strconv.Itoa(i.N)
	return fmt.Sprintf(":%s\r\n", nStr)
}

type BulkString struct {
	S *string
}

func (bs *BulkString) RespEncode() string {
	if bs.S == nil {
		return "$-1\r\n"
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(*bs.S), *bs.S)
}

type Array struct {
	Elements *[]RespEncoder
}

func (a *Array) RespEncode() string {
	if a.Elements == nil {
		return "*-1\r\n"
	}
	res := fmt.Sprintf("*%d\r\n", len(*a.Elements))
	for _, e := range *a.Elements {
		res = res + e.RespEncode()
	}
	return res
}

func (a *Array) AddElement(element RespEncoder) {
	if a.Elements == nil {
		a.Elements = new([]RespEncoder)
	}
	*a.Elements = append(*a.Elements, element)
}

func (a *Array) GetElements() []RespEncoder {
	if a.Elements == nil {
		return []RespEncoder{}
	}
	return *a.Elements
}

func RespDecode(s string) (RespEncoder, error) {
	switch {
	case strings.HasPrefix(s, "+"):
		return decodeSimpleString(s)
	case strings.HasPrefix(s, "-"):
		return decodeErrorString(s)
	case strings.HasPrefix(s, ":"):
		return decodeInteger(s)
	case strings.HasPrefix(s, "$"):
		return decodeBulkString(s)
	case strings.HasPrefix(s, "*"):
		return decodeArray(s)
	default:
		return nil, fmt.Errorf("invalid data type: unexpected prefix")
	}
}

func decodeSimpleString(s string) (*SimpleString, error) {
	s = strings.TrimPrefix(s, "+")
	splited := strings.Split(s, "\r\n")
	if len(splited) != 2 {
		return nil, fmt.Errorf("invalid SimpleString response: unexpected number of parts")
	}
	if splited[1] != "" {
		return nil, fmt.Errorf("invalid SimpleString response: unexpected string after split")
	}
	return &SimpleString{S: splited[0]}, nil
}

func decodeErrorString(s string) (*ErrorString, error) {
	s = strings.TrimPrefix(s, "-")
	splited := strings.Split(s, "\r\n")
	if len(splited) != 2 {
		return nil, fmt.Errorf("invalid ErrorString response: unexpected number of parts")
	}
	if splited[1] != "" {
		return nil, fmt.Errorf("invalid ErrorString response: unexpected string after split")
	}
	return &ErrorString{E: splited[0]}, nil
}

func decodeInteger(s string) (*Integer, error) {
	s = strings.TrimPrefix(s, ":")
	splited := strings.Split(s, "\r\n")
	if len(splited) != 2 {
		return nil, fmt.Errorf("invalid Integer response: unexpected number of parts")
	}
	if splited[1] != "" {
		return nil, fmt.Errorf("invalid Integer response: unexpected string after split")
	}
	n, err := strconv.Atoi(splited[0])
	if err != nil {
		return nil, fmt.Errorf("invalid Integer response: error converting string to integer: %w", err)
	}
	return &Integer{N: n}, nil
}

func decodeBulkString(s string) (*BulkString, error) {
	s = strings.TrimPrefix(s, "$")
	splited := strings.Split(s, "\r\n")
	if len(splited) == 2 {
		if splited[0] != "-1" || splited[1] != "" {
			return nil, fmt.Errorf("invalid null BulkString response: unexpected number of parts")
		}
		return &BulkString{}, nil
	}
	if len(splited) != 3 {
		return nil, fmt.Errorf("invalid BulkString response: unexpected number of parts")
	}
	l, err := strconv.Atoi(splited[0])
	if err != nil {
		return nil, fmt.Errorf("invalid BulkString response: error converting string to integer: %w", err)
	}
	if l != len(splited[1]) {
		return nil, fmt.Errorf("invalid BulkString response: mismatch of len of string: %d != %d", l, len(splited[1]))
	}
	return &BulkString{S: &splited[1]}, nil
}

func decodeArray(s string) (*Array, error) {
	s = strings.TrimPrefix(s, "*")
	splited := strings.Split(s, "\r\n")
	l, err := strconv.Atoi(splited[0])
	if err != nil {
		return nil, fmt.Errorf("invalid Array response: error converting string to integer: %w", err)
	}
	if l == -1 {
		return &Array{}, nil
	}
	return nil, fmt.Errorf("err")
}

