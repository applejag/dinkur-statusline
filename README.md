# Dinkur statusline

Small CLI to be used as a statusline in e.g:

- AwesomeWM (<https://awesomewm.org/>)

  - Example without coloring: <https://gitea.jillejr.tech/kalle/dotfiles/src/commit/dc34cd10837a9f2781fbe7c7377b14cde33fdc69/awesome/themes/holo/theme.lua#L353-L355>
  - Example using Pango coloring: <https://gitea.jillejr.tech/kalle/dotfiles/src/commit/8c57a1b3ef75d2d056848f19fa146ba810f75801/awesome/themes/holo/theme.lua#L353-L357>

- Pango Markup (<https://docs.gtk.org/Pango/pango_markup.html>), used in GTK apps and widgets, such as:

  - Executor - Gnome Shell Extension (<https://raujonas.github.io/executor/>)

  - Waybar - Wayland status bar (<https://github.com/Alexays/Waybar>)

    - Example using Pango coloring: <https://gitea.jillejr.tech/kalle/dotfiles/commit/1b47b12397a62267a38ac7fbeffcd2c3c3887a8d>

## Installation

Requires Go 1.18 (or higher):

```sh
go install github.com/dinkur/dinkur-statusline@latest
```

## Usage

```console
$ dinkur-statusline
1:01:32 (my task name) | 1:01:32 | 12%

$ dinkur-statusline --color pango
<span foreground='lime'>1:01:32 (my task name)</span> | 1:01:32 | 12%

$ dinkur-statusline --color raujonas-executor
<executor.markup.true> <span foreground='lime'>1:01:32 (my task name)</span> | 1:01:32 | 12%

$ dinkur-statusline --help
Usage: dinkur-statusline [flags]

Shows the status of your local Dinkur database (~/.local/share/dinkur/dinkur.db)
printed on a single line.

Possible colors:
  --color auto               Means "ansi" if interactive TTY, otherwise means "none"
  --color none               No coloring is applied
  --color pango              Coloring via Pango markup
  --color raujonas-executor  Same as "pango", but with added "<executor.markup.true>" prefix

Flags:
  -c, --color string      Color format (default "auto")
  -h, --help              Show this help text
      --work-hours uint   Hours in a workday, used in percentage calc (default 8)
```
