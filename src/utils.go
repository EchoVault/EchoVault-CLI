package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"

	"github.com/tidwall/resp"
)

func Contains[T comparable](arr []T, elem T) bool {
	for _, v := range arr {
		if v == elem {
			return true
		}
	}
	return false
}

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

	str += "\n"

	return str, nil
}

func Decode(raw string) ([]resp.Value, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	res := []resp.Value{}

	v, _, err := rd.ReadValue()

	if err != nil {
		return nil, err
	}

	if Contains[string]([]string{"SimpleString", "BulkString", "Integer", "Error"}, v.Type().String()) {
		return []resp.Value{v}, nil
	}

	if v.Type().String() == "Array" {
		for _, elem := range v.Array() {
			res = append(res, elem)
		}
	}

	return res, nil
}

func ReadMessage(r *bufio.ReadWriter) (message string, err error) {
	var line [][]byte

	for {
		b, _, err := r.ReadLine()

		if err != nil {
			return "", err
		}

		if bytes.Equal(b, []byte("")) {
			// End of message
			break
		}

		line = append(line, b)
	}

	return fmt.Sprintf("%s\r\n", string(bytes.Join(line, []byte("\r\n")))), nil
}

func PrintDecoded(arr []resp.Value, indent int) {
	if len(arr) == 0 {
		return
	}
	if len(arr) == 1 {
		if Contains([]string{"", "SUBSCRIBE_OK"}, arr[0].String()) {
			return
		}
		for i := 0; i < indent; i++ {
			fmt.Print(" ")
		}
		fmt.Println(arr[0])
		return
	}
	for _, item := range arr {
		if item.Type().String() == "Array" {
			PrintDecoded(item.Array(), indent+1)
			continue
		}
		if Contains[string]([]string{"SimpleString", "BulkString", "Integer", "Error"}, item.Type().String()) {
			PrintDecoded([]resp.Value{item}, indent)
		}
	}
}
