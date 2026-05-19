<div align="center">
  <h1>Zamaneh</h1>
  <h6>The way you manage your time</h6>
</div>
<div align="center">
  <img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/1995parham/zamaneh?sort=semver&style=for-the-badge&logo=github">
  <img alt="GitHub Release Date" src="https://img.shields.io/github/release-date/1995parham/zamaneh?style=for-the-badge&logo=github">
  <img alt="GitHub Actions Workflow Status" src="https://img.shields.io/github/actions/workflow/status/1995parham/zamaneh/go.yaml?style=for-the-badge&logo=github">
</div>

## Introduction

Zamaneh is a terminal pomodoro timer written in **Go 1.26**. It runs the
classic 25-min work / 5-min short break / 15-min long break cadence with a
beautiful Bubble Tea TUI, and appends each finished interval to a TOML log so
you can look back at how you spent your day.

This is the third life of the project: it started in Go, took a detour through
Rust, and is now back in Go with a richer TUI built on the Charm v2 stack
([Bubble Tea](https://charm.land/bubbletea/v2),
[Lip Gloss](https://charm.land/lipgloss/v2),
[Bubbles](https://charm.land/bubbles/v2)).

## Install

```sh
go install github.com/1995parham/zamaneh@latest
```

Or build from source:

```sh
go build -o zamaneh .
```

## Usage

```sh
zamaneh --title "writing zamaneh README"
```

Flags:

| Flag                 | Default | Description                                     |
| -------------------- | ------- | ----------------------------------------------- |
| `--title`            | `-`     | label for this working session                  |
| `--work`             | `25m`   | duration of a work interval                     |
| `--short-break`      | `5m`    | duration of a short break                       |
| `--long-break`       | `15m`   | duration of a long break                        |
| `--long-break-every` | `4`     | long break after this many work intervals       |
| `--verbose`          | `false` | enable debug logging                            |

While running:

| Key             | Action                  |
| --------------- | ----------------------- |
| `space` / `p`   | pause / resume          |
| `s`             | skip the current phase  |
| `r`             | reset the current phase |
| `t`             | rename the session      |
| `q` / `esc` / `ctrl+c` | save and quit    |

While renaming, `enter` saves the new title and `esc` cancels.

## Session log

Completed intervals are appended atomically to:

```
$XDG_DATA_HOME/zamaneh/sessions.toml
```

(falling back to `~/.local/share/zamaneh/sessions.toml` on Linux when the env
var is unset). The schema is intentionally simple:

```toml
[[sessions]]
title = "writing zamaneh README"
phase = "work"
started_at = 2026-05-19T10:30:00Z
ended_at = 2026-05-19T10:55:00Z
duration = 1500000000000
completed = true
```

## Acknowledgements

Built on the lovely [Charm](https://charm.land) tooling. XDG handling via
[`adrg/xdg`](https://github.com/adrg/xdg), TOML via
[`pelletier/go-toml`](https://github.com/pelletier/go-toml).
