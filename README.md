## EchoVault CLI Client

This is the CLI client for EchoVault. It allows you to interact
with the EchoVault server on your terminal.

You can view the EchoVault server repository [here](https://github.com/EchoVault/EchoVault).

## Installing

### Homebrew
1) `brew tap echovault/echovault`
2) `brew install echovault/echovault/echovault-cli`
3) `echovault-cli`

### Binaries
You can obtain the right binary for your system by clicking on a release
tag and scrolling to the assets section.

## Configuration

The EchoVault cli client takes the following args:

1) `--addr` - The IP address of the server. Default is `127.0.0.1`.
2) `--port` - The port to connect to on the server. Default is `7480`.
3) `--tls` - Boolean flag that instructs the client to establish a TLS connection with the server. Default is `false`.
4) `--mtls` - Boolean flag the instructs the client to establish an mTLS connection with the server. Default is `false`. If both `--tls` and `--mtls` are provided, `--mtls` will take priority.
5) `--server-ca` - The certificate authority used to verify the server on a TLS connection.
6) `--cert-key-pair` - Can specified multiple times. This is the comma-separated cert/key pair that the client will use to verify itself to the server on an mTLS connection. The format is `path/to/cert,path/to/key`