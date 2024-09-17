package shared

import (
	"os"

	"golang.org/x/exp/constraints"
	"golang.org/x/term"
)

type terminalOptions struct {
	MinWidth     int
	MaxWidth     int
	ClampWidth   bool
	DefaultWidth int
}

type TerminalOption func(*terminalOptions)

func TerminalDefaultWidth(size int) TerminalOption {
	return func(o *terminalOptions) {
		o.DefaultWidth = size
	}
}

func TerminalClampedWidth(min int, max int) TerminalOption {
	return func(o *terminalOptions) {
		o.MinWidth = min
		o.MaxWidth = max
		o.ClampWidth = true
	}
}

func GetTerminalWidth(options ...TerminalOption) int {
	opts := terminalOptions{
		DefaultWidth: 0,
		ClampWidth:   false,
	}

	for _, opt := range options {
		opt(&opts)
	}

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width = opts.DefaultWidth
	}

	if opts.ClampWidth {
		width = Clamp(width, opts.MinWidth, opts.MaxWidth)
	}

	return width
}

type Number interface {
	constraints.Integer | constraints.Float
}

func Min[T Number](value T, values ...T) T {
	length := len(values)
	if length < 1 {
		return value
	}

	for i := 0; i < length; i++ {
		if value > values[i] {
			value = values[i]
		}
	}

	return value
}

func Max[T Number](value T, values ...T) T {
	length := len(values)
	if length < 1 {
		return value
	}

	for i := 0; i < length; i++ {
		if value < values[i] {
			value = values[i]
		}
	}

	return value
}

func Clamp[T Number](value T, low T, high T) T {
	if low > high {
		high, low = low, high
	}

	return Min(Max(value, low), high)
}
