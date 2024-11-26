package config

import "time"

type Auth struct {
	JWT JWT `envPrefix:"JWT_"`
}

type JWT struct {
	Secret string `env:"SECRET,required"`

	LifetimeAccess  time.Duration `env:"LIFETIME_ACCESS,required"`
	LifetimeRefresh time.Duration `env:"LIFETIME_REFRESH,required"`
}
