# Binks CLI 🤖

**Codex‑style shell assistant written in Go**

> **Project status – Pre‑alpha.**  Expect rapid changes and rough edges while we stand up the very first commands.

Binks CLI is an experiment to recreate the OpenAI *Codex CLI* workflow using nothing but Go and open‑source tooling for a self-contained binary.

---

## 🚀 Quick start

```bash
# 1. Nightly binary (replace <platform>)
curl -L https://github.com/binks-cli/binks/releases/latest/download/binks_<platform>.tar.gz \
  | tar -xz -C /usr/local/bin

# 2. Or build from source
git clone https://github.com/binks-cli/binks && cd binks
go build -o binks ./cmd/binks
```

Launch an interactive session:

```bash
$ binks
binks:~/project >
```

---

## 🛠 Building & testing

```bash
go test ./...    # unit tests
go vet  ./...
```

*Requires Go 1.22 + and any POSIX‑style shell (Bash, Zsh, Fish, etc.).*

---

## 🧩 Architecture snapshot

```
┌───────────────────────┐
│        main.go        │
└──────────┬────────────┘
           │
    ┌──────▼─────┐   readline‑based
    │  REPL Loop │─────────────────────────────┐
    └──────┬─────┘                             │
  built‑ins │ external                         │
  ┌─────────▼───────────┐          ┌───────────▼───────────┐
  │   builtin cmds      │          │    executor.Exec       │
  │  (cd/exit/help)     │          │  wraps os/exec, Dir=cwd│
  └─────────────────────┘          └───────────────────────┘
```

- Single self‑contained Go binary
- Session object owns working directory & options
- Executor can run in *capture* (for tests) or *attached* (interactive programs) mode
- AI agent will plug into the REPL loop via a strategy interface

---

## 🗺 Roadmap

| Milestone                                           | Status         |
| --------------------------------------------------- | -------------- |
| **v0.1** – Minimal REPL & command execution         | 🚧 in progress |
| **v0.2** – Config file, prompt theming, first tests | 📝 planned     |
| **v0.3** – Safety heuristics & public preview       | 📝 planned     |
| **v1.0** – Optional AI agent (OpenAI / local)       | 🔜 back‑log    |

Live tracking: see the [Projects board](https://github.com/binks-cli/binks/projects).

---

## 🤝 Contributing

1. **Fork** → `git checkout -b feat/awesome`
2. **Write tests first** (`go test ./...`) – PRs without tests will be nudged 🙂
3. Follow [**Conventional Commits**](https://www.conventionalcommits.org/)
4. `go vet`, `go fmt`, make CI happy
5. Open a PR & describe *why* the change matters

Bug reports & feature ideas are very welcome. See [`CONTRIBUTING.md`](CONTRIBUTING.md) for full guidelines.

---

## 📜 License

Released under the **MIT License** – see [`LICENSE`](LICENSE).

> *Binks CLI is ****not**** affiliated with OpenAI or the original Codex CLI project.*\
> *“Codex” is a trademark of OpenAI; “Binks” is an independent, open‑source homage implemented in Go.*

---

### ⭐ If Binks helps you, star the repo and tell your friends!