package config

type Gateway struct {
	Minio Minio `envPrefix:"MINIO_"`
}

type Minio struct {
	Endpoint        string `env:"ENDPOINT,required"`
	AccessKeyID     string `env:"ACCESS_KEY_ID,required"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY,required"`
}
