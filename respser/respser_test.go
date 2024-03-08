package respser_test

import (
	"gored/respser"
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
