package respser

import (
	"fmt"
	"strconv"
)

func EncodeSimpleString(s string) string {
	return fmt.Sprintf("+%s\r\n", s)
}

func EncodeErrorString(s string) string {
	return fmt.Sprintf("-%s\r\n", s)
}

func EncodeInteger(n int) string {
	nStr := strconv.Itoa(n)
	return fmt.Sprintf(":%s\r\n", nStr)
}

func EncodeBulkStrings(s *string) string {
	if s == nil {
		return "$-1\r\n"
	}
	return fmt.Sprintf("$%d\r\n%s\r\n", len(*s), *s)
}
