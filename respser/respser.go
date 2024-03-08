package respser

import (
	"errors"
	"fmt"
	"strconv"
)

const CRLF = "\r\n"

var (
	ErrInvalidType       = errors.New("invalid type")
	ErrInvalidInputData  = errors.New("invalid input data")
	ErrDataMismatch      = errors.New("data mismatch")
	ErrInvalidInputParts = errors.New("invalid input parts")
	ErrFailedExtraction  = errors.New("failed extraction")
	ErrDecode            = errors.New("decode error")
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

func decodeError(fn string, in string) *RespSerError {
	return &RespSerError{fn, in, ErrDecode}
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
