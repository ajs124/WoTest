package config

type AuthenticationScheme uint

const (
	AuthNone = iota
	AuthBasic
	AuthDigest
	AuthBearer
	AuthApikey
	AuthPsk
	AuthOauth1
)

type AuthenticationData struct {
	Scheme AuthenticationScheme `json:"scheme"`
	Data   map[string]string    `json:"data"`
}
