package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Host string
	Port int
}

func Load() *Config {
	host := flag.String("host", getEnv("APP_HOST", "127.0.0.1"), "Server host")
	port := flag.Int("port", getEnvInt("APP_PORT", 8087), "Server port")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	return &Config{
		Host: *host,
		Port: *port,
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
