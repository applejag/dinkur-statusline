# Dinkur statusline

Small CLI to be used as a statusline in e.g:

- AwesomeWM (<https://awesomewm.org/>)

- Pango Markup (<https://docs.gtk.org/Pango/pango_markup.html>), used in GTK apps and widgets, such as:

  - Executor - Gnome Shell Extension (<https://raujonas.github.io/executor/>)

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
```
