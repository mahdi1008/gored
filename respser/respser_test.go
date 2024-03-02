package respser

import (
	"testing"
)

func TestEncodeSimpleString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{"normal_string", "Hello, World!", "+Hello, World!\r\n"},
		{"integer_string", "123", "+123\r\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := EncodeSimpleString(tc.input)
			if tc.want != res {
				t.Errorf("Expected: %s, Got: %s", tc.want, res)
			}
		})
	}
}

func TestEncodeErrorString(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  string
	}{
		{"normal_string", "Hello, World!", "-Hello, World!\r\n"},
		{"integer_string", "123", "-123\r\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := EncodeErrorString(tc.input)
			if tc.want != res {
				t.Errorf("Expected: %s, Got: %s", tc.want, res)
			}
		})
	}
}

func TestEncodeInteger(t *testing.T) {
	testCases := []struct {
		name  string
		input int
		want  string
	}{
		{"encode_positive_integer", 1000, ":1000\r\n"},
		{"encode_zero", 0, ":0\r\n"},
		{"encode_negative_integer", -3, ":-3\r\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := EncodeInteger(tc.input)
			if res != tc.want {
				t.Errorf("Expected: %s, Got: %s", tc.want, res)
			}
		})
	}
}

func TestEncodeBulkStrings(t *testing.T) {
	testCases := []struct {
		name       string
		bulkString string
		expectNil  bool
		want       string
	}{
		{"non-empty_string", "hello", false, "$5\r\nhello\r\n"},
		{"empty_string", "", false, "$0\r\n\r\n"},
		{"nil_string", "", true, "$-1\r\n"}, // Here we use an empty string to represent a nil bulkString.
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var bs *string
			if !tc.expectNil {
				bs = &tc.bulkString
			}
			res := EncodeBulkStrings(bs)
			if res != tc.want {
				t.Errorf("Expected: %s, Got: %s", tc.want, res)
			}
		})
	}
}
