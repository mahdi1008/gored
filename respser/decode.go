package respser

import (
	"errors"
	"strconv"
	"strings"
)

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
		return nil, invalidTypeError("RespDecode", s)
	}
}

func decodeSimpleString(s string) (*SimpleString, error) {
	s = strings.TrimPrefix(s, "+")
	splited := strings.Split(s, "\r\n")
	if len(splited) != 2 {
		return nil, invalidInputPartsError("decodeSimpleString", s)
	}
	if splited[1] != "" {
		return nil, invalidInputDataError("decodeSimpleString", s)
	}
	return &SimpleString{S: splited[0]}, nil
}

func decodeErrorString(s string) (*ErrorString, error) {
	s = strings.TrimPrefix(s, "-")
	splited := strings.Split(s, "\r\n")
	if len(splited) != 2 {
		return nil, invalidInputPartsError("decodeErrorString", s)
	}
	if splited[1] != "" {
		return nil, invalidInputDataError("decodeErrorString", s)
	}
	return &ErrorString{E: splited[0]}, nil
}

func decodeInteger(s string) (*Integer, error) {
	s = strings.TrimPrefix(s, ":")
	splited := strings.Split(s, "\r\n")
	if len(splited) != 2 {
		return nil, invalidInputPartsError("decodeInteger", s)
	}
	if splited[1] != "" {
		return nil, invalidInputDataError("decodeInteger", s)
	}
	n, err := strconv.Atoi(splited[0])
	if err != nil {
		return nil, errors.Join(invalidInputDataError("decodeInteger", s), err)
	}
	return &Integer{N: n}, nil
}

func decodeBulkString(s string) (*BulkString, error) {
	s = strings.TrimPrefix(s, "$")
	splited := strings.Split(s, "\r\n")
	if len(splited) == 2 {
		if splited[0] != "-1" || splited[1] != "" {
			return nil, invalidInputPartsError("decodeBulkString", s)
		}
		return &BulkString{}, nil
	}
	if len(splited) != 3 {
		return nil, invalidInputPartsError("decodeBulkString", s)
	}
	l, err := strconv.Atoi(splited[0])
	if err != nil {
		return nil, errors.Join(invalidInputDataError("decodeBulkString", s), err)
	}
	if l != len(splited[1]) {
		return nil, dataMismatchError("decodeBulkString", s)
	}
	return &BulkString{S: &splited[1]}, nil
}

func decodeArray(s string) (*Array, error) {
	splited := strings.Split(s, "\r\n")
	l, err := strconv.Atoi(strings.TrimPrefix(splited[0], "*"))
	if err != nil {
		return nil, errors.Join(invalidInputDataError("decodeArray", s), err)
	}
	if l == -1 {
		return &Array{}, nil
	}
	arr, r, err := ExtractArray(s)

	if err != nil {
		return nil, errors.Join(decodeError("decodeArray", s), err)
	}

	if r != "" {
		return nil, invalidInputDataError("decodeArray", s)
	}

	return arr, nil
}
