[![CI](https://github.com/binks-cli/binks/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/binks-cli/binks/actions/workflows/ci.yml)

# Binks CLI 

**A cross-platform, Go-powered re-imagining of Codename Goose & Codex CLI**

> **Project status – Alpha.**\
> The base architecture is implemented and stable. Expect some breaking changes as new features are added.

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
# Optional: static analysis (if installed)
golangci-lint run    
```

*Requires **Go 1.24+** and a POSIX-compatible shell (WSL / Git-Bash works on Windows).*

---

## 🧪 Continuous Integration (CI)

All tests are automatically run on every push and pull request via GitHub Actions. The CI workflow uses Go 1.24+ and runs `go test -v ./...` to ensure all tests pass in a clean environment.

To run the same tests locally:

```bash
go test ./...
```

If you add new tests, make sure they do not require manual input or special local setup. Use temporary files/directories and avoid prompts to ensure CI compatibility.

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
  built-ins│external                                       │
┌───────────▼───────────┐            ┌─────────────────────▼──────────────────-┐
│   cd / exit / help    │            │  executor.RunCmd (Stage1–3)             │
│   (session mutators)  │            │  - capture  (tests)                     │
│                       │            │  - attach   (interactive programs)      │
└───────────────────────┘            └─────────────────────────────────────────┘
                                         ▲
                                         │ later
                                         │
                                         │   ┌────────────── extension: shell / fs / git --┐
                                         └──▶│  AI Agent + MCP tool bus  (Stage6–7)        │
                                             └─────────────────────────────────────────────┘
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

## 🖥️ Non-Blocking GUI/Background Commands

When you launch a known GUI application or background process (e.g., `idea .`, `code .`, `chrome`, or `open` on macOS), Binks will start the process asynchronously and immediately return control to the prompt. You'll see a message like `[launched idea]` to confirm the launch. This prevents the shell from hanging while the app runs.

Other commands (including long-running ones like `sleep 10`) will block the prompt as usual. Support for explicit backgrounding with `&` is not yet implemented.

- If you run a command not in the known list, it will run synchronously by default.
- If you encounter a case where a GUI app blocks the prompt, please open an issue with details.

---

## Command History

Binks automatically saves your command history to a file in your home directory (`~/.binks_history`). This allows you to recall commands from previous sessions using the Up/Down arrows, similar to other shells. If you wish to clear your history, simply delete this file.

---
