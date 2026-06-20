# grout

A lightweight CLI for managing multi-service local dev environments via a single config file.

---

## Installation

```bash
go install github.com/yourname/grout@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/grout.git && cd grout && go build -o grout .
```

---

## Usage

Define your services in a `grout.yaml` file at the root of your project:

```yaml
services:
  api:
    cmd: "go run ./cmd/api"
    port: 8080
  worker:
    cmd: "go run ./cmd/worker"
    env:
      QUEUE_URL: "redis://localhost:6379"
  db:
    cmd: "docker run -p 5432:5432 postgres:15"
```

Then start all services with a single command:

```bash
grout up
```

**Common commands:**

| Command         | Description                        |
|-----------------|------------------------------------|
| `grout up`      | Start all services                 |
| `grout up api`  | Start a specific service           |
| `grout down`    | Stop all running services          |
| `grout status`  | Show status of all services        |
| `grout logs`    | Tail logs from all services        |

---

## License

MIT © [yourname](https://github.com/yourname)