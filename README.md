# Binks CLI 🦢

**A cross-platform, Go-powered re-imagining of Codename Goose & Codex CLI**

> **Project status – Pre-alpha.**\
> We are currently implementing the base architecture. Expect breaking changes.

Binks CLI lets you work inside a richer, single‑screen terminal UI that wraps **your existing shell** (bash, zsh, fish) across Linux, macOS, and Windows.\
It starts life as a **fast, self-contained shell wrapper** and will grow into a fully-featured AI agent with Model-Context-Protocol (MCP) extensions.

---

## 🚀 Quick start

```bash
# 1 · Download the latest nightly (replace <os>_<arch>)
curl -L https://github.com/binks-cli/binks/releases/latest/download/binks_<os>_<arch>.tar.gz \
  | tar -xz -C /usr/local/bin

# 2 · Or build from source
git clone https://github.com/binks-cli/binks && cd binks
go build -o binks ./cmd/binks
```

Launch an interactive session:

```bash
$ binks
binks:~/project >
```

-

---

## 🛠 Building & testing

```bash
go test ./...        # unit & integration tests 
go vet  ./...
golangci-lint run    # optional – static analysis
```

*Requires ****Go 1.22+**** and **``** in **``** (WSL / Git-Bash works on Windows).*

---

## 🧩 Architecture snapshot (MVP → future)

```
┌────────── main.go ──────────┐
│  flags / config bootstrap   │
└──────────┬──────────────────┘
           │
   ┌───────▼────────┐  readline / promptui
   │  REPL  (Stage2)│──────────────────────────────────────┐
   └───────┬────────┘                                      │
   built-ins│external                                      │
┌───────────▼───────────┐            ┌─────────────────────▼──────────────────┐
│   cd / exit / help    │            │  executor.RunCmd (Stage1–3)            │
│   (session mutators)  │            │  - capture  (tests)                    │
│                       │            │  - attach   (interactive programs)     │
└───────────────────────┘            └─────────────────────────────────────────┘
                                         ▲
                                         │ later
                                         │
                                         │   ┌────────────── extension: shell / fs / git … ┐
                                         └──▶│  AI Agent + MCP tool bus  (Stage6–7)        │
                                             └──────────────────────────────────────────────┘
```

- **Single static Go binary** – no Node/Rust runtime.
- **Session** owns working directory, history & options.
- **Executor** may *capture* output for tests or *attach* directly for editors/IDEs.
- **AI Agent** & **MCP extensions** plug in behind a strategy interface when enabled.

---

## 🗺 Roadmap & milestones

| Stage | Focus                                                     | Target tag  |
| ----- | --------------------------------------------------------- | ----------- |
| 1     | Basic command execution (`binks "echo hi"`)               | v0.0.1      |
| 2     | Interactive REPL loop + unit tests                        | v0.1.0      |
| 3     | Built-ins (`cd`, `exit`), session state, Windows support  | v0.2.0      |
| 4     | History, dynamic prompt, colour, interactive programs     | v0.3.0      |
| 5     | CI matrix (win/linux/mac), high test coverage             | v0.4.0      |
| 6     | LLM assistant (approval workflow, OpenAI chat)            | v0.5.0-beta |
| 7     | MCP extension bus, file-IO / git tools, multi-model agent | v1.0.0      |

> Detailed tasks live in the [GitHub Projects board](https://github.com/binks-cli/binks/projects).

---

## 🤝 Contributing

1. \*\*Fork → \*\*``
2. **Write tests first.** PRs without coverage will be asked to add it.
3. Follow **Conventional Commits**.
4. `go fmt`, `go vet`, ensure CI passes.
5. Open a PR describing **why** the change matters.

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the full guidelines.

---

## 📜 License

Released under **MIT**. See [`LICENSE`](LICENSE).

> Binks CLI is **not** affiliated with Block’s Codename Goose, Anthropic, or OpenAI.\
> “Binks” references the spirit of those projects; this is an independent re-implementation in Go.

---

### ⭐ If Binks CLI makes your day smoother, star the repo and share some love!
