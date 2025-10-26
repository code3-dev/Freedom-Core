# Freedom-Core 🚀

![Build](https://img.shields.io/github/actions/workflow/status/Freedom-Guard/freedom-core/go.yml?branch=main)
![Go Version](https://img.shields.io/badge/Go-1.22+-blue)
![License](https://img.shields.io/github/license/Freedom-Guard/Freedom-Core)
![Docker Pulls](https://img.shields.io/docker/pulls/freedom-guard/freedom-core)

**Freedom-Core** is an open-source, cross-platform tool designed to manage and execute multiple networking cores (like Xray, sing-box, or future modules) via a lightweight API and CLI. It provides a secure, flexible, and easy-to-use foundation for network operations, automation, and dashboard integration.

---

## 🌟 Key Features

- Manage multiple cores: Xray, sing-box, and future modules.
- Open a custom port to receive API commands.
- Cross-platform executable support: Windows (.exe), Linux, macOS.
- Lightweight Go-based service – can run natively or in containers.
- Optional Docker deployment for convenience.
- CLI interface for local control and notifications.
- Designed for integration with web dashboards or automation tools.

---

## 🛠️ Quick Start

### Prerequisites

- [Go 1.22+](https://golang.org/dl/) (for development and local builds)
- Optional: [Docker & Docker Compose](https://docs.docker.com/) (for containerized deployment)

---

### 2. Run Locally (Native)

```bash
go run cmd/server/main.go
```

Or build an executable:

```bash
# Windows 64-bit
go build -o freedom-core.exe ./cmd/server

# Windows 32-bit
GOOS=windows GOARCH=386 go build -o freedom-core-x86.exe ./cmd/server

# Linux 64-bit
go build -o freedom-core ./cmd/server

# Linux 32-bit
GOOS=linux GOARCH=386 go build -o freedom-core-linux-x86 ./cmd/server

# macOS Intel (x64)
GOOS=darwin GOARCH=amd64 go build -o freedom-core-macos-x64 ./cmd/server

# macOS Apple Silicon (ARM64)
GOOS=darwin GOARCH=arm64 go build -o freedom-core-macos-arm64 ./cmd/server
```

* The service will start and show:

```
Freedom-Core is running... 🚀
```

* By default, it opens port `8087` for API commands.

---

### 3. Run with Docker (Optional)

```bash
docker-compose up --build -d
```

* The service will run on port `8087`.
* Useful for isolated environments, but not required.

---

## 📌 Supported Cores

* **Xray** (vNext modules supported)
* **sing-box** (future)
* Other custom networking cores can be integrated via plugin modules.

---

## 📌 Roadmap

* [ ] Add CLI commands: `start`, `stop`, `status`
* [ ] Implement multi-core runner (Xray + sing-box + future)
* [ ] Full REST API: `/start`, `/stop`, `/config`, `/status`
* [ ] Web dashboard integration
* [ ] Publish first stable release (v0.1.0)
* [ ] Notifications for core events

---

## 📁 Project Structure

```
freedom-core/
├── cmd/                 # Entry point of the application
├── internal/            # Core logic (API, runner, notifier)
├── pkg/                 # Optional reusable packages
├── config/              # Configuration files (JSON/YAML)
├── docker/              # Docker-related files
├── .github/workflows/   # CI/CD configuration
├── Dockerfile
├── docker-compose.yml
├── README.md
└── go.mod
```

---

## 💡 Contribution

Freedom-Core is fully open-source! Contributions, ideas, and improvements are welcome.

* Fork the repository
* Create a new branch
* Commit your changes
* Submit a Pull Request

Please follow standard Go conventions and write clean, maintainable code.

---

## ⚖️ License

Freedom-Core is released under the [Apache 2.0](LICENSE). You are free to use, modify, and distribute it.

---

## 📣 Contact & Support

For questions, discussions, or feature requests:

* GitHub Issues: [https://github.com/Freedom-Guard/freedom-core/issues](https://github.com/Freedom-Guard/freedom-core/issues)
* GitHub Discussions (future): TBD