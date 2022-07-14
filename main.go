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
)

func main() {
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

	dayLength := time.Hour * 8
	dayPercentage := int64(100 * sumTimes / dayLength)

	var sb strings.Builder
	sb.WriteString("<executor.markup.true>")
	sb.WriteString(" <span foreground='")
	if activeEntry != nil {
		sb.WriteString("lime")
	} else {
		sb.WriteString("gray")
	}
	sb.WriteString("'>")
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
	sb.WriteString("</span>")

	sb.WriteString(" | <span>")
	sb.WriteString(FormatDuration(sumTimes))
	sb.WriteString("</span>")
	sb.WriteString(" | <span>")
	sb.WriteString(strconv.FormatInt(dayPercentage, 10))
	sb.WriteRune('%')
	sb.WriteString("</span>")

	if status.AFKSince != nil {
		if status.BackSince != nil {
			// has returned
			dur := status.BackSince.Sub(*status.AFKSince)
			sb.WriteString(" | <span foreground='orange'>AFK for ")
			sb.WriteString(FormatDuration(dur))
			sb.WriteString(" (welcome back)</span>")
		} else {
			// is AFK
			dur := time.Since(*status.AFKSince)
			sb.WriteString(" | <span foreground='red'>AFK for ")
			sb.WriteString(FormatDuration(dur))
			sb.WriteString("</span>")
		}
	}

	fmt.Print(sb.String())
}

func printErr(err error) {
	var sb strings.Builder
	sb.WriteString("<executor.markup.true>")
	sb.WriteString(" <span foreground='red'>err: ")
	sb.WriteString(err.Error())
	sb.WriteString("</span>")
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
