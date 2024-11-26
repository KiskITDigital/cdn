package config

type Transport struct {
	HTTP HTTP `envPrefix:"HTTP_"`
}

type HTTP struct {
	Port          int    `env:"PORT,required"`
	SwaggerUIPath string `env:"SWAGGER_UI_PATH"`
}
