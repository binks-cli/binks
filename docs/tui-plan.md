# TUI Plan and Alternate Screen Buffer (Issue #28)

## TUI Framework Research

- **Bubble Tea (Charmbracelet):**
  - Modern, flexible, and popular Go TUI framework.
  - State machine model, supports alternate screen, panes, and custom layouts.
  - Well-documented and actively maintained.
  - Best for future extensibility (AI panes, scrollback, etc).
- **tview (rivo):**
  - Widget-based, higher-level, easier for forms/tables.
  - Mature and widely used, but less flexible for custom layouts.

**Recommendation:**
Bubble Tea is the preferred choice for Binks due to its flexibility and suitability for custom, multi-pane UIs and future AI integration.

## Minimal Implementation (This Issue)

- Add support for alternate screen buffer mode, controlled by the `BINKS_ALT_SCREEN` environment variable.
- When enabled, Binks will:
  - Switch to the alternate screen buffer on startup (using ANSI code `\x1b[?1049h`).
  - Restore the normal screen buffer on exit (using ANSI code `\x1b[?1049l`).
- This is similar to how `less` or `vim` operate, and is a stepping stone to a full TUI.
- No new dependencies are required for this minimal step.

## High-Level Integration Plan

- In the future, refactor the main loop to use Bubble Tea’s model.
- Implement panes for scrollback/history, input, and AI output.
- Replace readline/prompt handling with Bubble Tea’s input system.
- Maintain an internal buffer for output, rendered by the TUI.

## Safety and Compatibility

- Alternate screen mode is opt-in and does not affect normal usage unless `BINKS_ALT_SCREEN=1` is set.
- If the terminal does not support ANSI codes, Binks will still function (but may display escape codes as text).
- No regressions to normal usage are expected.

## References
- [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- [tview](https://github.com/rivo/tview)

---

*This document summarizes the research and plan for full-screen TUI support in Binks, as outlined in [issue #28](https://github.com/binks-cli/binks/issues/28).*
