package respser_test

import (
	"gored/respser"
	"reflect"
	"testing"
)

func TestEncodeSimpleString(t *testing.T) {
	testCases := []struct {
		name  string
		input respser.SimpleString
		want  string
	}{
		{"normal_string", respser.SimpleString{S: "Hello, World!"}, "+Hello, World!\r\n"},
		{"integer_string", respser.SimpleString{S: "123"}, "+123\r\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.input.RespEncode()
			if tc.want != res {
				t.Errorf("Expected: %s, Got: %s", tc.want, res)
			}
		})
	}
}

func TestEncodeErrorString(t *testing.T) {
	testCases := []struct {
		name  string
		input respser.ErrorString
		want  string
	}{
		{"normal_string", respser.ErrorString{E: "Hello, World!"}, "-Hello, World!\r\n"},
		{"integer_string", respser.ErrorString{E: "123"}, "-123\r\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.input.RespEncode()
			if tc.want != res {
				t.Errorf("Expected: %s, Got: %s", tc.want, res)
			}
		})
	}
}

func TestEncodeInteger(t *testing.T) {
	testCases := []struct {
		name  string
		input respser.Integer
		want  string
	}{
		{"encode_positive_integer", respser.Integer{N: 1000}, ":1000\r\n"},
		{"encode_zero", respser.Integer{N: 0}, ":0\r\n"},
		{"encode_negative_integer", respser.Integer{N: -3}, ":-3\r\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.input.RespEncode()
			if res != tc.want {
				t.Errorf("Expected: %s, Got: %s", tc.want, res)
			}
		})
	}
}

func ptr(s string) *string {
	return &s
}

func TestEncodeBulkStrings(t *testing.T) {
	testCases := []struct {
		name  string
		input respser.BulkString
		want  string
	}{
		{"encode_non_empty_string", respser.BulkString{S: ptr("hello")}, "$5\r\nhello\r\n"},
		{"encode_empty_string", respser.BulkString{S: ptr("")}, "$0\r\n\r\n"},
		{"encode_nil_string", respser.BulkString{}, "$-1\r\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.input.RespEncode()
			if res != tc.want {
				t.Errorf("Expected: %s, Got: %s", tc.want, res)
			}
		})
	}
}

func TestEncodeArray(t *testing.T) {
	testCases := []struct {
		name  string
		input *respser.Array
		want  string
	}{
		{
			"encode_normal_array",
			&respser.Array{Elements: &[]respser.RespEncoder{&respser.Integer{N: 1}, &respser.Integer{N: 2}, &respser.Integer{N: 3}, &respser.Integer{N: 4}, &respser.BulkString{S: ptr("hello")}}},
			"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n",
		},
		{
			"encode_empty_array",
			&respser.Array{Elements: &[]respser.RespEncoder{}},
			"*0\r\n",
		},
		{
			"encode_nil_array",
			&respser.Array{},
			"*-1\r\n",
		},
		{
			"encode_nested_array",
			&respser.Array{Elements: &[]respser.RespEncoder{
				&respser.Array{
					Elements: &[]respser.RespEncoder{&respser.Integer{N: 1}, &respser.Integer{N: 2}, &respser.Integer{N: 3}},
				},
				&respser.Array{
					Elements: &[]respser.RespEncoder{&respser.SimpleString{S: "Hello"}, &respser.ErrorString{"World"}},
				},
			}},
			"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.input.RespEncode()
			if res != tc.want {
				t.Errorf("Expected: %s, Got: %s", tc.want, res)
			}
		})
	}
}

func TestDecodeSimpleString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  respser.SimpleString
	}{
		{"decode_normal_simple_string", "+hello world\r\n", respser.SimpleString{S: "hello world"}},
		{"decode_simple_string_with_empty_string", "+OK\r\n", respser.SimpleString{S: "OK"}},
		{"decode_simple_string_with_newline_in_string", "+hello\nworld\r\n", respser.SimpleString{S: "hello\nworld"}},
		{"decode_simple_string_with_carriage_return_in_string", "+hello\rworld\r\n", respser.SimpleString{S: "hello\rworld"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := respser.RespDecode(tc.input)
			if err != nil {
				t.Errorf("Expected nil err, Got %s", err)
			}

			if !reflect.DeepEqual(s, &tc.want) {
				t.Errorf("Expected SimpleString %v, Got %v", tc.want, s)
			}
		})
	}
}

func TestDecodeErrorString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  respser.ErrorString
	}{
		{"decode_normal_error_string", "-ERR unknown command 'foobar'\r\n", respser.ErrorString{E: "ERR unknown command 'foobar'"}},
		{"decode_error_string_with_empty_string", "-OK\r\n", respser.ErrorString{E: "OK"}},
		{"decode_error_string_with_newline_in_string", "-ERR hello\nworld\r\n", respser.ErrorString{E: "ERR hello\nworld"}},
		{"decode_error_string_with_carriage_return_in_string", "-ERR hello\rworld\r\n", respser.ErrorString{E: "ERR hello\rworld"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := respser.RespDecode(tc.input)
			if err != nil {
				t.Errorf("Expected nil err, Got %s", err)
			}

			if !reflect.DeepEqual(s, &tc.want) {
				t.Errorf("Expected ErrorString %v, Got %v", tc.want, s)
			}
		})
	}
}

func TestDecodeInteger(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  respser.Integer
	}{
		{"decode_normal_integer", ":1000\r\n", respser.Integer{N: 1000}},
		{"decode_integer_with_leading_zeros", ":00001000\r\n", respser.Integer{N: 1000}},
		{"decode_integer_with_negative_value", ":-1000\r\n", respser.Integer{N: -1000}},
		{"decode_integer_with_plus_sign", ":+1000\r\n", respser.Integer{N: 1000}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := respser.RespDecode(tc.input)
			if err != nil {
				t.Errorf("Expected nil err, Got %s", err)
			}

			if !reflect.DeepEqual(s, &tc.want) {
				t.Errorf("Expected Integer %v, Got %v", tc.want, s)
			}
		})
	}
}

func TestDecodeBulkString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  *respser.BulkString
	}{
		{"decode_normal_bulk_string", "$6\r\nfoobar\r\n", &respser.BulkString{S: ptr("foobar")}},
		{"decode_bulk_string_with_empty_string", "$0\r\n\r\n", &respser.BulkString{S: ptr("")}},
		{"decode_bulk_string_with_newline_in_string", "$7\r\nfoo\nbar\r\n", &respser.BulkString{S: ptr("foo\nbar")}},
		{"decode_bulk_string_with_carriage_return_in_string", "$7\r\nfoo\rbar\r\n", &respser.BulkString{S: ptr("foo\rbar")}},
		{"decode_null_bulk_string", "$-1\r\n", &respser.BulkString{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := respser.RespDecode(tc.input)
			if err != nil {
				t.Errorf("Expected nil err, Got %s", err)
			}

			if !reflect.DeepEqual(s, tc.want) {
				t.Errorf("Expected BulkString %v, Got %v", tc.want, s)
			}
		})
	}
}

// Decoder tests ---------

func TestRespDecodeErrorUnknownPrefix(t *testing.T) {
	input := "foo\r\n"
	if _, err := respser.RespDecode(input); err == nil {
		t.Errorf("Expected error for unknown prefix, got nil")
	}
}

func TestRespDecodeErrorInvalidSimpleStringParts(t *testing.T) {
	input := "+foo\r\nbar\r\n"
	if _, err := respser.RespDecode(input); err == nil {
		t.Errorf("Expected error for invalid simple string parts, got nil")
	}
}

func TestRespDecodeErrorInvalidErrorStringParts(t *testing.T) {
	input := "-foo\r\nbar\r\n"
	if _, err := respser.RespDecode(input); err == nil {
		t.Errorf("Expected error for invalid error string parts, got nil")
	}
}

func TestRespDecodeErrorInvalidIntegerParts(t *testing.T) {
	input := ":foo\r\n"
	if _, err := respser.RespDecode(input); err == nil {
		t.Errorf("Expected error for invalid integer parts, got nil")
	}
}

func TestRespDecodeErrorInvalidBulkStringParts(t *testing.T) {
	input := "$foo\r\nbar\r\n"
	if _, err := respser.RespDecode(input); err == nil {
		t.Errorf("Expected error for invalid bulk string parts, got nil")
	}
}

func TestRespDecodeErrorInvalidBulkStringLength(t *testing.T) {
	input := "$10\r\nfoo\r\n"
	if _, err := respser.RespDecode(input); err == nil {
		t.Errorf("Expected error for invalid bulk string length, got nil")
	}
}

func TestRespDecodeErrorInvalidBulkStringNull(t *testing.T) {
	input := "$0\r\nfoo\r\n"
	if _, err := respser.RespDecode(input); err == nil {
		t.Errorf("Expected error for invalid bulk string null, got nil")
	}
}
