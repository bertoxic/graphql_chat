package config

type AppConfig struct {
	DataBaseURL string
	JWTSecret   string
	Port        string
}

func NewConfig(databaseURL, jwtSecret, port string) AppConfig {
	return AppConfig{
		DataBaseURL: databaseURL,
		JWTSecret:   jwtSecret,
		Port:        port,
	}
}
