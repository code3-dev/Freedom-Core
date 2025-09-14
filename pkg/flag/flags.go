package flags

import (
	"flag"
	"fmt"
	"os"

	"github.com/Freedom-Guard/freedom-core/pkg/updater"
)

type Config struct {
	Host string
	Port int
}

var Version = "0.2.0"

func Parse() *Config {
	host := flag.String("host", "127.0.0.1", "Server host")
	port := flag.Int("port", 8087, "Server port")
	showVersion := flag.Bool("version", false, "Show Freedom Core version")
	update := flag.Bool("update", false, "Update Freedom Core to the latest version")

	flag.Parse()

	if *showVersion {
		fmt.Println("Freedom Core by Freedom Guard")
		fmt.Printf("Version: %s\n", Version)
		fmt.Println("GitHub: https://github.com/Freedom-Guard/Freedom-Core")
		os.Exit(0)
	}

	if *update {
		fmt.Println("Checking for updates...")
		if err := updater.Update(); err != nil {
			fmt.Println("Update failed:", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	return &Config{
		Host: *host,
		Port: *port,
	}
}
