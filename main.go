package main

import (
	"context"
	"fmt"
	"html"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dinkur/dinkur/pkg/dinkur"
	"github.com/dinkur/dinkur/pkg/dinkurdb"
	"github.com/dinkur/dinkur/pkg/timeutil"
	"github.com/mattn/go-colorable"
	"github.com/spf13/pflag"
)

var flags = struct {
	color     string
	showHelp  bool
	workHours uint
}{
	workHours: 8,
}

type colorType int

const (
	colorNone colorType = iota
	colorAnsi
	colorPango
	colorRaujonasExecutor
)

var coloring colorType

func init() {
	pflag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: dinkur-statusline [flags]

Shows the status of your local Dinkur database (~/.local/share/dinkur/dinkur.db)
printed on a single line.

Possible colors:
  --color auto               Means "ansi" if interactive TTY, otherwise means "none"
  --color none               No coloring is applied
  --color pango              Coloring via Pango markup
  --color raujonas-executor  Same as "pango", but with added "<executor.markup.true>" prefix

Flags:
`)
		pflag.PrintDefaults()
	}
	pflag.StringVarP(&flags.color, "color", "c", "auto", "Color format")
	pflag.BoolVarP(&flags.showHelp, "help", "h", false, "Show this help text")
	pflag.UintVar(&flags.workHours, "work-hours", flags.workHours, "Hours in a workday, used in percentage calc")
}

func main() {
	pflag.Parse()
	if flags.showHelp {
		pflag.Usage()
		os.Exit(0)
	}

	switch strings.ToLower(strings.TrimSpace(flags.color)) {
	case "auto":
		if noColor {
			coloring = colorNone
		} else {
			coloring = colorAnsi
		}
	case "none":
		coloring = colorNone
	case "ansi":
		coloring = colorAnsi
	case "pango":
		coloring = colorPango
	case "raujonas-executor":
		coloring = colorRaujonasExecutor
	default:
		fmt.Fprintf(os.Stderr, "error: unknown coloring type: %q", flags.color)
		fmt.Fprint(os.Stderr, `Possible values:
  --color auto               Means "ansi" if interactive TTY, otherwise means "none"
  --color none               No coloring is applied
  --color pango              Coloring via Pango markup
  --color raujonas-executor  Same as "pango", but with added "<executor.markup.true>" prefix
`)
		os.Exit(1)
	}

	defer func() {
		if p := recover(); p != nil {
			printErr(fmt.Errorf("panic: %v", p))
			os.Exit(0)
		}
	}()

	home, err := os.UserHomeDir()
	if err != nil {
		printErr(err)
		os.Exit(1)
	}
	opt := dinkurdb.Options{
		MkdirAll:             false,
		SkipMigrateOnConnect: true,
		DebugLogging:         false,
	}
	db := dinkurdb.NewClient(filepath.Join(home, ".local", "share", "dinkur", "dinkur.db"), opt)

	ctx := context.Background()

	if err := db.Connect(ctx); err != nil {
		printErr(err)
		os.Exit(1)
	}

	tasks, err := db.GetEntryList(ctx, dinkur.SearchEntry{
		Shorthand: timeutil.TimeSpanThisDay,
	})
	if err != nil {
		printErr(err)
		os.Exit(1)
	}
	status, err := db.GetStatus(ctx)
	if err != nil {
		printErr(err)
		os.Exit(1)
	}

	var activeEntry *dinkur.Entry
	var sumTimes time.Duration
	for _, task := range tasks {
		if task.End == nil {
			activeEntry = &task
		}
		sumTimes += task.Elapsed()
	}

	dayLength := time.Hour * time.Duration(flags.workHours)
	dayPercentage := int64(100 * sumTimes / dayLength)

	var sb strings.Builder
	if coloring == colorRaujonasExecutor {
		sb.WriteString("<executor.markup.true> ")
	}
	switch coloring {
	case colorAnsi:
		if activeEntry != nil {
			sb.WriteString(ansiFgGreen)
		} else {
			sb.WriteString(ansiFgHiBlack)
		}
	case colorPango, colorRaujonasExecutor:
		sb.WriteString("<span foreground='")
		if activeEntry != nil {
			sb.WriteString("lime")
		} else {
			sb.WriteString("gray")
		}
		sb.WriteString("'>")
	}
	if activeEntry != nil {
		sb.WriteString(FormatDuration(activeEntry.Elapsed()))
		if activeEntry.Name != "" {
			sb.WriteString(" (")
			sb.WriteString(html.EscapeString(activeEntry.Name))
			sb.WriteRune(')')
		}
	} else {
		sb.WriteString("no active task")
	}

	switch coloring {
	case colorAnsi:
		sb.WriteString(ansiReset)
	case colorPango, colorRaujonasExecutor:
		sb.WriteString("</span>")
	}

	sb.WriteString(" | ")
	sb.WriteString(FormatDuration(sumTimes))
	sb.WriteString(" | ")
	sb.WriteString(strconv.FormatInt(dayPercentage, 10))
	sb.WriteRune('%')

	if status.AFKSince != nil {
		if status.BackSince != nil {
			// has returned
			dur := status.BackSince.Sub(*status.AFKSince)

			sb.WriteString(" | ")
			switch coloring {
			case colorAnsi:
				sb.WriteString(ansiFgYellow + "AFK for ")
				sb.WriteString(FormatDuration(dur))
				sb.WriteString(" (welcome back)" + ansiReset)
			case colorPango, colorRaujonasExecutor:
				sb.WriteString("<span foreground='orange'>AFK for ")
				sb.WriteString(FormatDuration(dur))
				sb.WriteString(" (welcome back)</span>")
			}
		} else {
			// is AFK
			dur := time.Since(*status.AFKSince)
			switch coloring {
			case colorAnsi:
				sb.WriteString(" | " + ansiFgHiRed + "AFK for ")
				sb.WriteString(FormatDuration(dur))
				sb.WriteString(ansiReset)
			case colorPango, colorRaujonasExecutor:
				sb.WriteString(" | <span foreground='red'>AFK for ")
				sb.WriteString(FormatDuration(dur))
				sb.WriteString("</span>")
			}
		}
	}

	if coloring == colorAnsi {
		fmt.Fprintln(colorable.NewColorableStdout(), sb.String())
	} else {
		fmt.Println(sb.String())
	}
}

func printErr(err error) {
	switch coloring {
	case colorAnsi:
		fmt.Fprintf(colorable.NewColorableStderr(), "%serr: %s%s", ansiFgHiRed, err, ansiReset)
	case colorPango, colorRaujonasExecutor:
		fmt.Fprintf(os.Stderr, "<executor.markup.true> <span foreground='red'>err: %s</span>", err)
	default:
		fmt.Fprintf(os.Stderr, "err: %s", err)
	}
}

type cmdResult struct {
	stderr string
	stdout string
}

func execCmd(name string, args ...string) (cmdResult, error) {
	cmd := exec.Command(name, args...)

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return cmdResult{}, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return cmdResult{}, err
	}

	if err := cmd.Start(); err != nil {
		return cmdResult{}, err
	}

	stderr, err := io.ReadAll(stderrPipe)
	if err != nil {
		return cmdResult{}, err
	}

	stdout, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return cmdResult{}, err
	}

	if err := cmd.Wait(); err != nil {
		return cmdResult{}, err
	}

	return cmdResult{
		stdout: string(stdout),
		stderr: string(stderr),
	}, nil
}

func FormatDuration(d time.Duration) string {
	totalSeconds := int64(d / time.Second)
	var b []byte

	b = strconv.AppendInt(b, totalSeconds/3600, 10)

	b = append(b, ':')

	minutes := totalSeconds / 60 % 60
	if minutes < 10 {
		b = append(b, '0')
	}
	b = strconv.AppendInt(b, minutes, 10)

	b = append(b, ':')

	seconds := totalSeconds % 60
	if seconds < 10 {
		b = append(b, '0')
	}
	b = strconv.AppendInt(b, seconds, 10)

	return string(b)
}
