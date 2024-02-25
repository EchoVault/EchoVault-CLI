run:
	go run ./src/... --addr=localhost --port=7479

run-tls:
	go run ./src/... --addr=localhost --port=7479 --server-ca=./openssl/server/rootCA.crt --tls

run-mtls:
	go run ./src/... --addr=localhost --port=7479 --cert-key-pair=./openssl/client/client1.crt,./openssl/client/client1.key --server-ca=./openssl/server/rootCA.crt --tls --mtls
