package config

type Service struct {
	Salt          string `env:"FMT_SALT" envDefault:"MyUniqueSalt"`
	AdminUsername string `env:"DB_USERNAME" envDefault:"Admin"`
	AdminPassword string `env:"DB_PASS" envDefault:"qwerty"`
	AdminEmail    string `env:"DB_Email" envDefault:"qwerty@email.com"`
}
