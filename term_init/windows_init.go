//go:build windows

package term_init

import (
	"os"

	"golang.org/x/sys/windows"
)

func EnableTerm() {
	stdout := windows.Handle(os.Stdout.Fd())
	stdin := windows.Handle(os.Stdin.Fd())

	var originalMode uint32
	var originalModeIn uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)

	windows.GetConsoleMode(stdin, &originalModeIn)
	if (originalModeIn & (1 << 9)) != 0 {
		windows.SetConsoleMode(stdin, originalModeIn^windows.ENABLE_VIRTUAL_TERMINAL_INPUT)
	}
}

func StartInput() {
	stdin := windows.Handle(os.Stdin.Fd())

	var originalModeIn uint32
	windows.SetConsoleMode(stdin, originalModeIn|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	windows.SetConsoleMode(stdin, originalModeIn^windows.ENABLE_VIRTUAL_TERMINAL_INPUT)
}

func ResetTerm() {
	stdin := windows.Handle(os.Stdin.Fd())
	var originalModeIn uint32
	windows.GetConsoleMode(stdin, &originalModeIn)
	if (originalModeIn & (1 << 9)) != 0 {
		windows.SetConsoleMode(stdin, originalModeIn^windows.ENABLE_VIRTUAL_TERMINAL_INPUT)
	}
	os.Stdout.WriteString("\033[?25h")
}
