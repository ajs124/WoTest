package main

type AuthenticationScheme int8

const (
	None = iota
	Basic
	Digest
	Bearer
	Apikey
	Psk
	Oauth1
)

type AuthenticationData struct {
	scheme AuthenticationScheme
	// maybe that should be map[string][]byte ?
	data map[string]string
}
