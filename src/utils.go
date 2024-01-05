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

	str += "\r\n"

	return str, nil
}

func Decode(raw string) (resp.Value, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))
	var res resp.Value

	v, _, err := rd.ReadValue()

	if err != nil {
		return resp.Value{}, err
	}

	if Contains[string]([]string{"SimpleString", "BulkString", "Integer", "Error"}, v.Type().String()) {
		return v, nil
	}

	if v.Type().String() == "Array" {
		res = v
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
