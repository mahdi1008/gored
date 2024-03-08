package respser_test

import (
	"errors"
	"gored/respser"
	"reflect"
	"testing"
)

func TestRespDecode(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  respser.RespEncoder
		err   error
	}{
		{"decode_invalid_type", "hello world\r\n", nil, respser.ErrInvalidType},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := respser.RespDecode(tc.input)

			if !reflect.DeepEqual(s, tc.want) {
				t.Errorf("Expected %v, Got %v", tc.want, s)
			}

			if !errors.Is(err, tc.err) {
				t.Errorf("Expected error %v, Got %v", tc.err, err)
			}
		})
	}
}

func TestDecodeSimpleString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *respser.SimpleString
		err   error
	}{
		{"decode_normal_simple_string", "+hello world\r\n", &respser.SimpleString{S: "hello world"}, nil},
		{"decode_simple_string_with_empty_string", "+OK\r\n", &respser.SimpleString{S: "OK"}, nil},
		{"decode_simple_string_with_newline_in_string", "+hello\nworld\r\n", &respser.SimpleString{S: "hello\nworld"}, nil},
		{"decode_simple_string_with_carriage_return_in_string", "+hello\rworld\r\n", &respser.SimpleString{S: "hello\rworld"}, nil},
		{"decode_invalid_simple_string_with_extra_parts", "+hello\r\nworld\r\n", nil, respser.ErrInvalidInputParts},
		{"decode_invalid_simple_string_with_data_after_split", "+hello\r\nworld", nil, respser.ErrInvalidInputData},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := respser.RespDecode(tc.input)

			if !reflect.DeepEqual(s, tc.want) {
				t.Errorf("Expected SimpleString %v, Got %v", tc.want, s)
			}

			if !errors.Is(err, tc.err) {
				t.Errorf("Expected error %v, Got %v", tc.err, err)
			}
		})
	}
}

func TestDecodeErrorString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *respser.ErrorString
		err   error
	}{
		{"decode_normal_error_string", "-ERR unknown command 'foobar'\r\n", &respser.ErrorString{E: "ERR unknown command 'foobar'"}, nil},
		{"decode_error_string_with_empty_string", "-OK\r\n", &respser.ErrorString{E: "OK"}, nil},
		{"decode_error_string_with_newline_in_string", "-ERR hello\nworld\r\n", &respser.ErrorString{E: "ERR hello\nworld"}, nil},
		{"decode_error_string_with_carriage_return_in_string", "-ERR hello\rworld\r\n", &respser.ErrorString{E: "ERR hello\rworld"}, nil},
		{"decode_invalid_error_string_with_extra_parts", "-ERR hello\r\nworld\r\n", nil, respser.ErrInvalidInputParts},
		{"decode_invalid_error_string_with_data_after_split", "-ERR hello\r\nworld", nil, respser.ErrInvalidInputData},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := respser.RespDecode(tc.input)

			if !reflect.DeepEqual(s, tc.want) {
				t.Errorf("Expected ErrorString %v, Got %v", tc.want, s)
			}

			if !errors.Is(err, tc.err) {
				t.Errorf("Expected error %v, Got %v", tc.err, err)
			}
		})
	}
}

func TestDecodeInteger(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *respser.Integer
		err   error
	}{
		{"decode_normal_integer", ":1000\r\n", &respser.Integer{N: 1000}, nil},
		{"decode_integer_with_leading_zeros", ":00001000\r\n", &respser.Integer{N: 1000}, nil},
		{"decode_integer_with_negative_value", ":-1000\r\n", &respser.Integer{N: -1000}, nil},
		{"decode_integer_with_plus_sign", ":+1000\r\n", &respser.Integer{N: 1000}, nil},
		{"decode_invalid_parts", ":4\r\n3\r\n", nil, respser.ErrInvalidInputParts},
		{"decode_invalid_data_after_crlf", ":4\r\n3", nil, respser.ErrInvalidInputData},
		{"decode_string_instead_of_integer", ":foo\r\n", nil, respser.ErrInvalidInputData},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := respser.RespDecode(tc.input)

			if !reflect.DeepEqual(s, tc.want) {
				t.Errorf("Expected Integer %v, Got %v", tc.want, s)
			}

			if !errors.Is(err, tc.err) {
				t.Errorf("Expected error %v, Got %v", tc.err, err)
			}
		})
	}
}

func TestRespDecodeBulkString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *respser.BulkString
		err   error
	}{
		{"decode_normal_bulk_string", "$6\r\nfoobar\r\n", &respser.BulkString{S: ptr("foobar")}, nil},
		{"decode_bulk_string_with_empty_string", "$0\r\n\r\n", &respser.BulkString{S: ptr("")}, nil},
		{"decode_null_bulk_string", "$-1\r\n", &respser.BulkString{}, nil},
		{"decode_bulk_string_with_newline_in_string", "$7\r\nfoo\nbar\r\n", &respser.BulkString{S: ptr("foo\nbar")}, nil},
		{"decode_bulk_string_with_carriage_return_in_string", "$7\r\nfoo\rbar\r\n", &respser.BulkString{S: ptr("foo\rbar")}, nil},
		{"decode_error_invalid_bulk_string_parts", "$foo\r\nbar\r\n", nil, respser.ErrInvalidInputData},
		{"decode_error_invalid_bulk_string_length", "$10\r\nfoo\r\n", nil, respser.ErrDataMismatch},
		{"decode_error_invalid_bulk_string_null", "$0\r\nfoo\r\nsalam\r\n", nil, respser.ErrInvalidInputParts},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := respser.RespDecode(tc.input)
			if err == nil && !reflect.DeepEqual(s, tc.want) {
				t.Errorf("Expected BulkString %v, Got %v", tc.want, s)
			}

			if !errors.Is(err, tc.err) {
				t.Errorf("Expected error %v, Got %v", tc.err, err)
			}
		})
	}
}
func TestDecodeArray(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *respser.Array
		err   error
	}{
		{
			"decode_normal_array",
			"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			&respser.Array{
				Elements: &[]respser.RespEncoder{
					&respser.BulkString{S: ptr("foo")}, &respser.BulkString{S: ptr("bar")},
				},
			},
			nil,
		},
		{
			"decode_array_with_empty_item",
			"*1\r\n$0\r\n\r\n",
			&respser.Array{
				Elements: &[]respser.RespEncoder{
					&respser.BulkString{S: ptr("")},
				},
			},
			nil,
		},
		{
			"decode_array_with_null_item",
			"*1\r\n$-1\r\n",
			&respser.Array{
				Elements: &[]respser.RespEncoder{
					&respser.BulkString{},
				},
			},
			nil,
		},
		{
			"decode_array_with_null_item_and_empty_item",
			"*2\r\n$-1\r\n$0\r\n\r\n",
			&respser.Array{
				Elements: &[]respser.RespEncoder{
					&respser.BulkString{},
					&respser.BulkString{S: ptr("")},
				},
			},
			nil,
		},
		{
			"decode_array_with_nested_array",
			"*2\r\n*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n*2\r\n$3\r\nbaz\r\n$3\r\nqux\r\n",
			&respser.Array{
				Elements: &[]respser.RespEncoder{
					&respser.Array{
						Elements: &[]respser.RespEncoder{
							&respser.BulkString{S: ptr("foo")},
							&respser.BulkString{S: ptr("bar")},
						},
					},
					&respser.Array{
						Elements: &[]respser.RespEncoder{
							&respser.BulkString{S: ptr("baz")},
							&respser.BulkString{S: ptr("qux")},
						},
					},
				},
			}, nil},
		{
			"decode_array_with_newline_in_item",
			"*2\r\n$7\r\nfoo\nbar\r\n$3\r\nqux\r\n",
			&respser.Array{
				Elements: &[]respser.RespEncoder{
					&respser.BulkString{S: ptr("foo\nbar")},
					&respser.BulkString{S: ptr("qux")},
				},
			}, nil},
		{
			"decode_array_with_carriage_return_in_item",
			"*2\r\n$7\r\nfoo\rbar\r\n$3\r\nqux\r\n",
			&respser.Array{
				Elements: &[]respser.RespEncoder{
					&respser.BulkString{S: ptr("foo\rbar")},
					&respser.BulkString{S: ptr("qux")},
				},
			}, nil},
		{
			"decode_error_invalid_array_parts",
			"*foo\r\nbar\r\n",
			nil,
			respser.ErrInvalidInputData,
		},
		{
			"decode_error_invalid_array_length",
			"*10\r\n$3\r\nfoo\r\n",
			nil,
			respser.ErrFailedExtraction,
		},
		{
			"decode_error_invalid_array_extra_input",
			"*2\r\n$3\r\nfoo\r\n$3\r\nfoo\r\n\n",
			nil,
			respser.ErrInvalidInputData,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := respser.RespDecode(tc.input)
			if err == nil && !reflect.DeepEqual(s, tc.want) {
				t.Errorf("Expected Array %v, Got %v", tc.want, s)
			}

			if !errors.Is(err, tc.err) {
				t.Errorf("Expected error %v, Got %v", tc.err, err)
			}
		})
	}
}
