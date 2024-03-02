package respser

import (
	"testing"
)

func TestEncodeSimpleString(t *testing.T) {
	testCases := []struct {
		name  string
		input SimpleString
		want  string
	}{
		{"normal_string", SimpleString{S: "Hello, World!"}, "+Hello, World!\r\n"},
		{"integer_string", SimpleString{S: "123"}, "+123\r\n"},
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
		input ErrorString
		want  string
	}{
		{"normal_string", ErrorString{E: "Hello, World!"}, "-Hello, World!\r\n"},
		{"integer_string", ErrorString{E: "123"}, "-123\r\n"},
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
		input Integer
		want  string
	}{
		{"encode_positive_integer", Integer{N: 1000}, ":1000\r\n"},
		{"encode_zero", Integer{N: 0}, ":0\r\n"},
		{"encode_negative_integer", Integer{N: -3}, ":-3\r\n"},
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

func stringPointer(s string) *string {
	return &s
}

func TestEncodeBulkStrings(t *testing.T) {
	testCases := []struct {
		name  string
		input BulkString
		want  string
	}{
		{"encode_non_empty_string", BulkString{S: stringPointer("hello")}, "$5\r\nhello\r\n"},
		{"encode_empty_string", BulkString{S: stringPointer("")}, "$0\r\n\r\n"},
		{"encode_nil_string", BulkString{}, "$-1\r\n"},
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
		input *Array
		want  string
	}{
		{
			"encode_normal_array",
			&Array{Elements: &[]RespEncoder{&Integer{N: 1}, &Integer{N: 2}, &Integer{N: 3}, &Integer{N: 4}, &BulkString{S: stringPointer("hello")}}},
			"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n",
		},
		{
			"encode_empty_array",
			&Array{Elements: &[]RespEncoder{}},
			"*0\r\n",
		},
		{
			"encode_nil_array",
			&Array{},
			"*-1\r\n",
		},
		{
			"encode_nested_array",
			&Array{Elements: &[]RespEncoder{
				&Array{
					Elements: &[]RespEncoder{&Integer{N: 1}, &Integer{N: 2}, &Integer{N: 3}},
				},
				&Array{
					Elements: &[]RespEncoder{&SimpleString{S: "Hello"}, &ErrorString{"World"}},
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
