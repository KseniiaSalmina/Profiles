package config

type Formatter struct {
	Salt string `env:"FMT_SALT" envDefault:"MyUniqueSalt"`
}
