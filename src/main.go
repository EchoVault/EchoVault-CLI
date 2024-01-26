package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	conf := GetConfig()

	var conn net.Conn
	var err error

	// Writers & readers for stdio
	stdout, stdin := io.Writer(os.Stdin), io.Reader(os.Stdout)

	dialer := net.Dialer{
		KeepAlive: 200 * time.Millisecond,
	}

	if !conf.TLS {
		stdout.Write([]byte("Establishing TCP connection...\n"))
		conn, err = dialer.Dial("tcp", fmt.Sprintf("%s:%d", conf.Addr, conf.Port))
		if err != nil {
			panic(err)
		}
	} else {
		// Dial TLS
		stdout.Write([]byte("Establishing TLS connection...\n"))

		f, err := os.Open(conf.Cert)
		if err != nil {
			panic(err)
		}

		cert, err := io.ReadAll(bufio.NewReader(f))
		if err != nil {
			panic(err)
		}

		rootCAs := x509.NewCertPool()
		ok := rootCAs.AppendCertsFromPEM(cert)
		if !ok {
			panic("Failed to parse certificate")
		}

		conn, err = tls.DialWithDialer(
			&dialer,
			"tcp",
			fmt.Sprintf("%s:%d", conf.Addr, conf.Port),
			&tls.Config{
				RootCAs: rootCAs,
			})

		if err != nil {
			panic(fmt.Sprintf("Handshake Error: %s", err.Error()))
		}
	}

	defer conn.Close()

	done := make(chan struct{})

	// Writers & readers for connection
	cw, cr := io.Writer(conn), io.Reader(conn)

	go func() {
		for {
			stdout.Write([]byte("\n> "))

			if in, err := ReadMessage(stdin, []byte{'\n'}); err != nil {
				log.Println(err)
			} else {
				in := bytes.TrimSpace(in)

				if bytes.Equal(bytes.ToLower(in), []byte("quit\n\x00\x00\x00")) {
					break
				}

				// Serialize command and send to connection
				encoded, err := Encode(string(in))

				if err != nil {
					fmt.Println(err)
					continue
				}

				if _, err := cw.Write([]byte(encoded)); err != nil {
					log.Println(err)
				}

				// Read response from server
				message, err := ReadMessage(cr, []byte{'\r', '\n', '\r', '\n'})

				if err != nil && err == io.EOF {
					log.Println("connection closed")
					break
				} else if err != nil {
					log.Println(err)
				}

				decoded, err := Decode(message)
				if err != nil {
					log.Println(err)
					continue
				}

				if IsSubscribeResponse(decoded) {
					// If we're subscribed to a channel, listen for messages from the channel
					func() {
						for {
							var message []byte

							if msg, err := ReadMessage(cr, []byte{'\r', '\n', '\r', '\n'}); err != nil {
								if err == io.EOF {
									return
								}
								log.Println(err)
								continue
							} else {
								message = msg
							}

							decoded, err := Decode(message)
							if err != nil {
								log.Println(err)
								continue
							}

							cw.Write([]byte("+ACK\r\n\r\n"))
							if !decoded.IsNull() {
								PrintDecoded(decoded)
							}
						}
					}()
				} else {
					PrintDecoded(decoded)
				}
			}
		}
		done <- struct{}{}
	}()

	<-done
}
