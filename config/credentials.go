// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Credentials struct {
	AccessKey string `config:"accessKey"`
	SecretKey string `config:"secretKey"`
	Token     string `config:"token"`
}
