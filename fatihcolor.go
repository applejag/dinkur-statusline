// The MIT License (MIT)
//
// Copyright (c) 2013 Fatih Arslan
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"os"

	"github.com/mattn/go-isatty"
)

// NoColor defines if the output is colorized or not. It's dynamically set to
// false or true based on the stdout's file descriptor referring to a terminal
// or not. It's also set to true if the NO_COLOR environment variable is
// set (regardless of its value). This is a global option and affects all
// colors. For more control over each color block use the methods
// DisableColor() individually.
var noColor = noColorExists() || os.Getenv("TERM") == "dumb" ||
	(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))

func noColorExists() bool {
	_, exists := os.LookupEnv("NO_COLOR")
	return exists
}

const (
	ansiReset     = "\x1b[0m"

	ansiFgBlack   = "\x1b[30m"
	ansiFgRed     = "\x1b[31m"
	ansiFgGreen   = "\x1b[32m"
	ansiFgYellow  = "\x1b[33m"
	ansiFgBlue    = "\x1b[34m"
	ansiFgMagenta = "\x1b[35m"
	ansiFgCyan    = "\x1b[36m"
	ansiFgWhite   = "\x1b[37m"

	ansiFgHiBlack   = "\x1b[90m"
	ansiFgHiRed     = "\x1b[91m"
	ansiFgHiGreen   = "\x1b[92m"
	ansiFgHiYellow  = "\x1b[93m"
	ansiFgHiBlue    = "\x1b[94m"
	ansiFgHiMagenta = "\x1b[95m"
	ansiFgHiCyan    = "\x1b[96m"
	ansiFgHiWhite   = "\x1b[97m"
)
