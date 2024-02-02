package main

import (
	"encoding/json"
	"errors"
	"flag"
	yaml "gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
	"strings"
)

type Config struct {
	TLS          bool       `json:"TLS" yaml:"TLS"`
	MTLS         bool       `json:"MTLS" yaml:"MTLS"`
	CertKeyPairs [][]string `json:"CertKeyPairs" yaml:"CertKeyPairs"`
	ServerCAs    []string   `json:"ServerCAs" yaml:"ServerCAs"`
	Port         uint16     `json:"Port" yaml:"Port"`
	Addr         string     `json:"Addr" yaml:"Addr"`
}

func GetConfig() Config {
	var certKeyPairs [][]string
	var serverCAs []string

	flag.Func("certKeyPair",
		"A cert/key pair used by the server to verify the client. The value is 2 comma separated file paths.",
		func(s string) error {
			pair := strings.Split(strings.TrimSpace(s), ",")
			for i := 0; i < len(pair); i++ {
				pair[i] = strings.TrimSpace(pair[i])
			}
			if len(pair) != 2 {
				return errors.New("certKeyPair must be 2 comma separated file paths")
			}
			certKeyPairs = append(certKeyPairs, pair)
			return nil
		})

	flag.Func("serverCA",
		"A file path to a root CA used by the client to verify the server.",
		func(s string) error {
			serverCAs = append(serverCAs, s)
			return nil
		})

	tls := flag.Bool("tls", false, "Start the server in TLS mode. Default is false.")
	mtls := flag.Bool("mtls", false, "Use mTLS to verify the client with the server.")
	port := flag.Int("port", 7480, "Port to use. Default is 7480.")
	config := flag.String(
		"config",
		"",
		`File path to a JSON or YAML config file.The values in this config file will override the flag values.`,
	)
	addr := flag.String("addr", "127.0.0.1", "On src, this is the address of a server node to connect to.")

	flag.Parse()

	var conf Config

	if len(*config) > 0 {
		// Load config from config file
		f, err := os.Open(*config)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err = f.Close(); err != nil {
				log.Println(err)
			}
		}()

		ext := path.Ext(f.Name())

		if ext == ".json" {
			json.NewDecoder(f).Decode(&conf)
		}

		if ext == ".yaml" || ext == ".yml" {
			yaml.NewDecoder(f).Decode(&conf)
		}

		return conf
	}

	conf = Config{
		CertKeyPairs: certKeyPairs,
		ServerCAs:    serverCAs,
		TLS:          *tls,
		MTLS:         *mtls,
		Addr:         *addr,
		Port:         uint16(*port),
	}

	return conf
}
