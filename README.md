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

---

## ⚙️ Configuration & Environment Variables

Binks supports the following environment variables:

| Variable            | Purpose                                                      | Default                  |
|---------------------|--------------------------------------------------------------|--------------------------|
| `OPENAI_API_KEY`    | Use a real OpenAI agent for AI mode.                         | (unset = stub agent)     |
| `OPENAI_MODEL`      | Model name for OpenAI integration.                           | `gpt-3.5-turbo`          |
| `OPENAI_API_BASE`   | Override OpenAI API base URL.                                | `https://api.openai.com/v1` |
| `BINKS_ALT_SCREEN`  | Set to `1` to enable alternate screen buffer (TUI prep).     | `0` (disabled)           |
| `BINKS_DEBUG_AI`    | Set to `1` for debug logs from the AI agent.                 | `0` (disabled)           |

See `env.example` for a template.

---

## 🖥️ Alternate Screen / TUI Plan

Binks supports an alternate screen buffer (like `vim`/`less`) when `BINKS_ALT_SCREEN=1` is set. This is a stepping stone to a full TUI (see [TUI Plan](docs/tui-plan.md)).

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

1. **Fork → clone**
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

**Current async commands:** `idea`, `code`, `chrome`, `open`

Other commands (including long-running ones like `sleep 10`) will block the prompt as usual. Support for explicit backgrounding with `&` is not yet implemented.

- If you run a command not in the known list, it will run synchronously by default.
- If you encounter a case where a GUI app blocks the prompt, please open an issue with details.

---

## 🖥️ Interactive Program Support

Binks supports running many interactive and full-screen console programs (like `vim`, `nano`, `less`, `man`, `ssh`, `top`, etc.) directly from the REPL. When you run one of these commands, Binks will yield full control of the terminal to the program until it exits, then restore your prompt and session state. This means you can use editors, pagers, and SSH sessions as you would in a normal shell.

**Supported interactive commands include:**
- vim, nvim, vi
- nano
- less, more, man
- ssh
- top, htop

If you find an interactive program that does not work as expected, please open an issue. For best results, ensure your terminal supports ANSI escape codes and is not running in a restricted environment.

---

## Command History

Binks automatically saves your command history to a file in your home directory (`~/.binks_history`). This allows you to recall commands from previous sessions using the Up/Down arrows, similar to other shells. If you wish to clear your history, simply delete this file.

---

## 🤖 AI Mode

Binks supports an **AI mode** that lets you route your input to an AI agent for natural language queries, code suggestions, or shell command help.

### How to use AI mode

- **Per-query:** Prefix your input with `>>` to send it to the AI agent. Example:
  ```
  >> How do I list all files including hidden ones?
  ```
  The agent will respond with an answer or suggestion.

- **Global AI mode:** Enable AI mode for all input with the command:
  ```
  :ai on
  ```
  When enabled, the prompt will show `[AI]` and all input will be sent to the agent unless you prefix with `!` to force a shell command (e.g., `!ls`).
  To turn off AI mode, use:
  ```
  :ai off
  ```

- **Prompt indication:**
  - When AI mode is active, the prompt changes to `[AI] binks:~/dir >` (with color if your terminal supports it).
  - When AI mode is off, the prompt is `binks:~/dir >`.

- **AI agent:**
  - By default, Binks uses a stub agent that echoes your query. If you set `OPENAI_API_KEY`, Binks will use the real OpenAI API for responses.

- **Error handling:**
  - If the agent is unavailable or returns an error, you’ll see a clear error message.

### AI Command Suggestion Confirmation (v0.5.0+)

When the AI agent responds with a shell command suggestion (in a code block), Binks will **never execute it automatically**. Instead, you will see:

- The AI's explanation (if present)
- The suggested command, clearly formatted
- A confirmation prompt: `Execute this? [y/N]:`

**Example:**
```
[AI] Here is what you should do:
AI suggests: git pull && make test
Execute this? [y/N]:
```

- Type `y` or `yes` to approve and run the command.
- Type `n`, `no`, or just press Enter to decline (the command will not run).
- If you decline, you'll see `[AI] Cancelled.`

This workflow ensures you are always in control of what gets executed, even when using powerful AI agents.

- Only the first code block in the AI response is considered a command suggestion.
- If no code block is present, the AI's response is shown as plain text.
- Declined suggestions are not logged by default (see roadmap for future enhancements).

---
