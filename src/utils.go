package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/tidwall/resp"
)

func tokenize(comm string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(comm))
	r.Comma = ' '
	return r.Read()
}

func Encode(comm string) (string, error) {
	tokens, err := tokenize(comm)

	if err != nil {
		return "", errors.New("could not parse command")
	}

	str := fmt.Sprintf("*%d\r\n", len(tokens))

	for i, token := range tokens {
		if i == 0 {
			str += fmt.Sprintf("$%d\r\n%s\r\n", len(token), strings.ToUpper(token))
		} else {
			str += fmt.Sprintf("$%d\r\n%s\r\n", len(token), token)
		}
	}

	str += "\r\n"

	return str, nil
}

func Decode(raw []byte) (resp.Value, error) {
	rd := resp.NewReader(bytes.NewBuffer(raw))
	var res resp.Value

	v, _, err := rd.ReadValue()

	if err != nil {
		return resp.Value{}, err
	}

	if slices.Contains([]string{"SimpleString", "BulkString", "Integer", "Error"}, v.Type().String()) {
		return v, nil
	}

	if v.Type().String() == "Array" {
		res = v
	}

	return res, nil
}

func ReadMessage(r io.Reader, delim []byte) ([]byte, error) {
	buffSize := 8
	buff := make([]byte, buffSize)

	var n int
	var err error
	var res []byte

	for {
		n, err = r.Read(buff)
		res = append(res, buff...)
		if n < buffSize || err != nil {
			break
		}
		if bytes.Equal(buff[len(buff)-len(delim):], delim) {
			break
		}
		clear(buff)
	}

	return res, err
}

func IsSubscribeResponse(val resp.Value) bool {
	if val.Type().String() != "SimpleString" {
		return false
	}
	return val.String() == "SUBSCRIBE_OK"
}

func PrintArray(val resp.Value, initialIndent int) {
	if len(val.Array()) == 0 {
		fmt.Println("(empty array)")
		return
	}
	for i, item := range val.Array() {
		if i > 0 {
			// Prepend initial indent
			for j := 0; j < initialIndent; j++ {
				fmt.Print(" ")
			}
		}
		pos := fmt.Sprintf("%d) ", i+1)
		fmt.Print(pos)
		if item.Type().String() == "Array" {
			PrintArray(item, initialIndent+len(pos))
			continue
		}
		fmt.Println(item)
	}
}

func PrintDecoded(val resp.Value) {
	switch val.Type().String() {
	default:
		if val.IsNull() {
			fmt.Println("(nil)")
			return
		}
		fmt.Println(val)
	case "Array":
		PrintArray(val, 0)
	}
}
