package flags

type Config struct {
    Host string
    Port int
    Version string
}

var AppConfig *Config
