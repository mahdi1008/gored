package respser

import (
	"fmt"
	"strconv"
)

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
