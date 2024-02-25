run:
	go run ./src/... --addr=localhost --port=7479

run-tls:
	go run ./src/... --addr=localhost --port=7479 --serverCA=./openssl/server/rootCA.crt --tls

run-mtls:
	go run ./src/... --addr=localhost --port=7479 --certKeyPair=./openssl/client/client1.crt,./openssl/client/client1.key --serverCA=./openssl/server/rootCA.crt --tls --mtls
