package config

type Service struct {
	Salt string `env:"FMT_SALT" envDefault:"MyUniqueSalt"`
}
