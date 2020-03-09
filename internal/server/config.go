package server

type Config struct {
	BindAddr       string `env:"BIND_ADDR" envDefault:":8080"`
	DownloadedPath string `env:"DOWNLOADED_PATH" envDefault:"./downloaded/"`
}

func NewConfig() *Config {
	return &Config{}
}
