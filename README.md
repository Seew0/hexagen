# hexagen — Golang Hexagonal Architecture Generator

```
 _   _                                 
| | | | _____  ____ _  __ _  ___ _ __  
| |_| |/ _ \ \/ / _` |/ _` |/ _ \ '_ \ 
|  _  |  __/>  < (_| | (_| |  __/ | | |
|_| |_|\___/_/\_\__,_|\__, |\___|_| |_|
                      |___/            
```


`hexagen` is a CLI tool that scaffolds a production-ready **Golang Hexagonal Architecture** service using:

- Gin (HTTP router)
- Uber FX (dependency injection)
- Zap (logging)
- Config injection
- Go-embed templates

It generates a clean microservice skeleton with best practices built-in.

---

## 🚀 Installation

Install globally:

```
go install github.com/seew0/hexagen@latest
```

Check version:

```
hexagen --version
```

---

## 🧰 Usage

Generate a new service:

```
hexagen -r myservice -m github.com/me/myservice -p 8080
```

Interactive mode:

```
hexagen -i
```

Show version:

```
hexagen --version
```

---

## 🎛 CLI Flags

| Flag | Description |
|------|-------------|
| `-r` | Target directory |
| `-m` | Module name |
| `-p` | Server port |
| `-g` | Add `.gitkeep` |
| `-c` | Clean directory |
| `-i` | Interactive mode |
| `--version` | Show version |

---

## 📁 Generated structure

```
myservice/
└── cmd/
    └── main.go
└── commons/
    └── utils/
        └── logger.go
└── config/
    └── init/
        └── serverConfig.go
└── services/
    └── serviceName/
        └── routes/
            └── router.go
└── templates/
└── go.mod
└── Makefile
```

---

## 🧩 What's included

- Gin router
- Uber FX DI setup
- Lifecycle hooks
- Zap logger provider
- Config provider (ENV, SERVICE_NAME, PORT)
- Routing module
- Makefile + go.mod setup
- Embedded templates

---

## 📄 Template system

All templates live under:

```
templates/
- main.go.tmpl
- router.go.tmpl
- serverConfig.go.tmpl
- logger.go.tmpl
```

They are embedded using Go’s `embed.FS`.

---

## 🧪 Generated endpoints

```
GET /
→ { "status": "ok" }

GET /api/v1/ping
→ { "status": "ok", "pong": true }
```

---

## 🔌 Plugin system (kubectl-style)

Any executable named:

```
hexagen-<pluginname>
```

automatically becomes a subcommand:

```
hexagen <pluginname>
```

### Example

Create plugin:

```
hexagen-module
```

Make it executable:

```
chmod +x hexagen-module
```

Plugin contents:

```
#!/bin/bash
NAME=$1
mkdir -p services/$NAME
echo "Created module: $NAME"
```

Run:

```
hexagen module user
```

This executes `hexagen-module user`.

## 🤝 Contributing

1. Modify templates in `/templates`
2. Extend the generator logic
3. Add plugin binaries (`hexagen-xxx`)
4. Submit PRs

---

## 📜 License

MIT