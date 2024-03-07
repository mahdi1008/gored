package respser_test

import (
	"errors"
	"gored/respser"
	"reflect"
	"testing"
)

func TestExtractSimpleString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *respser.SimpleString
		want1 string
		err   error
	}{
		{"extract_normal_simple_string", "+This is a simple string\r\n", &respser.SimpleString{S: "This is a simple string"}, "", nil},
		{"extract_simple_string_with_empty_string", "+\r\n", &respser.SimpleString{S: ""}, "", nil},
		{"extract_simple_string_with_newline_in_string", "+This is a simple string\r\nAnd this is the remainder\r\n", &respser.SimpleString{S: "This is a simple string"}, "And this is the remainder\r\n", nil},
		{"extract_simple_string_with_carriage_return_in_string", "+This is a simple string\r\nAnd this \r is the remainder", &respser.SimpleString{S: "This is a simple string"}, "And this \r is the remainder", nil},
		{"extract_invalid_type_simple_string", "This is not a simple string", nil, "This is not a simple string", respser.ErrInvalidType},
		{"extract_invalid_input_parts_simple_string", "+This is not a simple string", nil, "+This is not a simple string", respser.ErrInvalidInputParts},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, got1, got2 := respser.ExtractSimpleString(tc.input)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Expected SimpleString %v, Got %v", tc.want, got)
			}

			if got1 != tc.want1 {
				t.Errorf("Expected remainder %v, Got %v", tc.want1, got1)
			}

			if !errors.Is(got2, tc.err) {
				t.Errorf("Expected error %v, Got %v", tc.err, got2)
			}
		})
	}
}

func TestExtractErrorString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *respser.ErrorString
		want1 string
		err   error
	}{
		{"extract_normal_error_string", "-This is an error string\r\n", &respser.ErrorString{E: "This is an error string"}, "", nil},
		{"extract_error_string_with_empty_string", "-\r\n", &respser.ErrorString{E: ""}, "", nil},
		{"extract_error_string_with_newline_in_string", "-This is an error string\r\nAnd this is the remainder\r\n", &respser.ErrorString{E: "This is an error string"}, "And this is the remainder\r\n", nil},
		{"extract_error_string_with_carriage_return_in_string", "-This is an error string\r\nAnd this \r is the remainder", &respser.ErrorString{E: "This is an error string"}, "And this \r is the remainder", nil},
		{"extract_invalid_type_error_string", "This is not an error string", nil, "This is not an error string", respser.ErrInvalidType},
		{"extract_invalid_input_parts_error_string", "-This is not an error string", nil, "-This is not an error string", respser.ErrInvalidInputParts},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, got1, got2 := respser.ExtractErrorString(tc.input)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Expected ErrorString %v, Got %v", tc.want, got)
			}

			if got1 != tc.want1 {
				t.Errorf("Expected remainder %v, Got %v", tc.want1, got1)
			}

			if !errors.Is(got2, tc.err) {
				t.Errorf("Expected error %v, Got %v", tc.err, got2)
			}
		})
	}
}

func TestExtractInteger(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *respser.Integer
		want1 string
		err   error
	}{
		{"extract_normal_integer", ":12345\r\n", &respser.Integer{N: 12345}, "", nil},
		{"extract_invalid_integer_with_empty_string", ":\r\n", nil, ":\r\n", respser.ErrInvalidInputData},
		{"extract_integer_with_newline_in_string", ":12345\r\nAnd this is the remainder\r\n", &respser.Integer{N: 12345}, "And this is the remainder\r\n", nil},
		{"extract_integer_with_carriage_return_in_string", ":12345\r\nAnd this \r is the remainder", &respser.Integer{N: 12345}, "And this \r is the remainder", nil},
		{"extract_invalid_type_integer", "This is not an integer", nil, "This is not an integer", respser.ErrInvalidType},
		{"extract_invalid_input_parts_integer", ":12345", nil, ":12345", respser.ErrInvalidInputParts},
		{"extract_invalid_input_type_integer", ":123.45\r\n", nil, ":123.45\r\n", respser.ErrInvalidInputData},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, got1, got2 := respser.ExtractInteger(tc.input)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Expected Integer %v, Got %v", tc.want, got)
			}

			if got1 != tc.want1 {
				t.Errorf("Expected remainder %v, Got %v", tc.want1, got1)
			}

			if !errors.Is(got2, tc.err) {
				t.Errorf("Expected error %v, Got %v", tc.err, got2)
			}
		})
	}
}

func TestExtractBulkString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *respser.BulkString
		want1 string
		err   error
	}{
		{"extract_normal_bulk_string", "$6\r\nfoobar\r\n", &respser.BulkString{S: ptr("foobar")}, "", nil},
		{"extract_normal_bulk_string_with_extra", "$6\r\nfoobar\r\n+hello\r\n-world\r\n", &respser.BulkString{S: ptr("foobar")}, "+hello\r\n-world\r\n", nil},
		{"extract_bulk_string_with_empty_string", "$0\r\n\r\n", &respser.BulkString{S: ptr("")}, "", nil},
		{"extract_bulk_string_with_newline_in_string", "$7\r\nfoo\nbar\r\n", &respser.BulkString{S: ptr("foo\nbar")}, "", nil},
		{"extract_bulk_string_with_carriage_return_in_string", "$7\r\nfoo\rbar\r\n", &respser.BulkString{S: ptr("foo\rbar")}, "", nil},
		{"extract_invalid_type_bulk_string", "This is not a bulk string", nil, "This is not a bulk string", respser.ErrInvalidType},
		{"extract_invalid_input_parts_bulk_string", "$6", nil, "$6", respser.ErrInvalidInputParts},
		{"extract_invalid_input_type_bulk_string", "$-1.5\r\nhi\r\n", nil, "$-1.5\r\nhi\r\n", respser.ErrInvalidInputData},
		{"extract_data_mismatch_bulk_string", "$6\r\nfoo\r\n", nil, "$6\r\nfoo\r\n", respser.ErrDataMismatch},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, got1, got2 := respser.ExtractBulkString(tc.input)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Expected BulkString %v, Got %v", tc.want, got)
			}

			if got1 != tc.want1 {
				t.Errorf("Expected remainder %v, Got %v", tc.want1, got1)
			}

			if !errors.Is(got2, tc.err) {
				t.Errorf("Expected error %v, Got %v", tc.err, got2)
			}
		})
	}
}

func TestExtractArray(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *respser.Array
		want1 string
		err   error
	}{
		{"extract_normal_array", "*3\r\n:1\r\n:2\r\n:3\r\n", &respser.Array{Elements: &[]respser.RespEncoder{&respser.Integer{N: 1}, &respser.Integer{N: 2}, &respser.Integer{N: 3}}}, "", nil},
		{"extract_array_with_empty_elements", "*0\r\n", &respser.Array{}, "", nil},
		{"extract_array_with_different_types_in_elements", "*4\r\n:1\r\n+hello!\r\n-world!\r\n$2\r\nhi\r\n", &respser.Array{Elements: &[]respser.RespEncoder{&respser.Integer{N: 1}, &respser.SimpleString{S: "hello!"}, &respser.ErrorString{E: "world!"}, &respser.BulkString{S: ptr("hi")}}}, "", nil},
		{"extract_array_with_recursive_array_in_elements", "*3\r\n*3\r\n:1\r\n:2\r\n:3\r\n:2\r\n:3\r\n", &respser.Array{Elements: &[]respser.RespEncoder{&respser.Array{Elements: &[]respser.RespEncoder{&respser.Integer{N: 1}, &respser.Integer{N: 2}, &respser.Integer{N: 3}}}, &respser.Integer{N: 2}, &respser.Integer{N: 3}}}, "", nil},
		{"extract_invalid_type_array", "This is not an array", nil, "This is not an array", respser.ErrInvalidType},
		{"extract_invalid_input_parts_array", "*3", nil, "*3", respser.ErrInvalidInputParts},
		{"extract_array_from_invalid_input", "*3\r\n:1\r\n:2\r\n:3\r\n:4", &respser.Array{Elements: &[]respser.RespEncoder{&respser.Integer{N: 1}, &respser.Integer{N: 2}, &respser.Integer{N: 3}}}, ":4", nil},
		{"extract_failed_extraction_array", "*3\r\n:1\r\n:2\r\nThis is not an integer\r\n:4", nil, "*3\r\n:1\r\n:2\r\nThis is not an integer\r\n:4", respser.ErrFailedExtraction},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, got1, got2 := respser.ExtractArray(tc.input)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Expected Array %v, Got %v", tc.want, got)
			}

			if got1 != tc.want1 {
				t.Errorf("Expected remainder %v, Got %v", tc.want1, got1)
			}

			if !errors.Is(got2, tc.err) {
				t.Errorf("Expected error %v, Got %v", tc.err, got2)
			}
		})
	}
}
