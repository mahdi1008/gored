package respser

import (
	"errors"
	"strconv"
	"strings"
)

func ExtractType(s string) (RespEncoder, string, error) {
	switch {
	case strings.HasPrefix(s, "+"):
		return ExtractSimpleString(s)
	case strings.HasPrefix(s, "-"):
		return ExtractErrorString(s)
	case strings.HasPrefix(s, ":"):
		return ExtractInteger(s)
	case strings.HasPrefix(s, "$"):
		return ExtractBulkString(s)
	case strings.HasPrefix(s, "*"):
		return ExtractArray(s)
	default:
		return nil, s, invalidTypeError("extractType", s)
	}
}

func ExtractSimpleString(s string) (*SimpleString, string, error) {
	if !strings.HasPrefix(s, "+") {
		return nil, s, invalidTypeError("extractSimpleString", s)
	}
	splited := strings.SplitN(s, CRLF, 2)
	if len(splited) != 2 {
		return nil, s, invalidInputPartsError("extractSimpleString", s)
	}
	internalString := strings.TrimPrefix(splited[0], "+")
	ss := &SimpleString{
		S: internalString,
	}
	return ss, splited[1], nil
}

func ExtractErrorString(s string) (*ErrorString, string, error) {
	if !strings.HasPrefix(s, "-") {
		return nil, s, invalidTypeError("extractErrorString", s)
	}
	splited := strings.SplitN(s, CRLF, 2)
	if len(splited) != 2 {
		return nil, s, invalidInputPartsError("extractErrorString", s)
	}
	internalString := strings.TrimPrefix(splited[0], "-")
	ss := &ErrorString{
		E: internalString,
	}
	return ss, splited[1], nil
}

func ExtractInteger(s string) (*Integer, string, error) {
	if !strings.HasPrefix(s, ":") {
		return nil, s, invalidTypeError("extractInteger", s)
	}
	splited := strings.SplitN(s, CRLF, 2)
	if len(splited) != 2 {
		return nil, s, invalidInputPartsError("extractInteger", s)
	}
	internalIntString := strings.TrimPrefix(splited[0], ":")
	internalInteger, err := strconv.Atoi(internalIntString)
	if err != nil {
		return nil, s, invalidInputDataError("extractInteger", s)
	}
	ss := &Integer{
		N: internalInteger,
	}
	return ss, splited[1], nil
}

func ExtractBulkString(s string) (*BulkString, string, error) {
	if !strings.HasPrefix(s, "$") {
		return nil, s, invalidTypeError("extractBulkString", s)
	}
	splited := strings.SplitN(s, CRLF, 3)
	if len(splited) == 2 {
		if splited[0] != "$-1" || splited[1] != "" {
			return nil, s, invalidInputPartsError("decodeBulkString", s)
		}
		return &BulkString{}, splited[1], nil
	}
	if len(splited) != 3 {
		return nil, s, invalidInputPartsError("extractBulkString", s)
	}
	internalBulkStringSizeString := strings.TrimPrefix(splited[0], "$")
	internalBulkStringSize, err := strconv.Atoi(internalBulkStringSizeString)
	if err != nil {
		return nil, s, invalidInputDataError("extractBulkString", s)
	}
	if internalBulkStringSize == -1 {
		return &BulkString{S: nil}, splited[1] + CRLF + splited[2], nil
	}

	internalString := splited[1]
	if len(internalString) != internalBulkStringSize {
		return nil, s, dataMismatchError("extractBulkString", s)
	}
	ss := &BulkString{
		S: &internalString,
	}
	return ss, splited[2], nil
}

func ExtractArray(s string) (*Array, string, error) {
	if !strings.HasPrefix(s, "*") {
		return nil, s, invalidTypeError("extractArray", s)
	}
	splited := strings.SplitN(s, CRLF, 2)
	if len(splited) != 2 {
		return nil, s, invalidInputPartsError("extractArray", s)
	}
	internalSizeString := strings.TrimPrefix(splited[0], "*")
	internalSize, err := strconv.Atoi(internalSizeString)
	if err != nil {
		return nil, s, invalidInputDataError("extractArray", s)
	}

	a := &Array{}
	remainder := splited[1]

	for i := 0; i < internalSize; i++ {
		re, r, err := ExtractType(remainder)
		if err != nil {
			return nil, s, errors.Join(failedExtractionError("extractArray", s), err)
		}
		a.AddElement(re)
		remainder = r
	}

	return a, remainder, nil
}
