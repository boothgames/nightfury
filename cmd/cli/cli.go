package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var whiteColor *color.Color
var redColor *color.Color
var yellowColor *color.Color
var greenColor *color.Color

func init() {
	whiteColor = color.New(color.FgWhite)
	redColor = color.New(color.FgRed)
	yellowColor = color.New(color.FgYellow)
	greenColor = color.New(color.FgGreen)
}

// Success print values to os.Stdout in green color
func Success(value ...interface{}) {
	stream := os.Stdout
	_, _ = greenColor.Add(color.Bold).Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
}

// Info print values to os.Stdout in white color
func Info(value ...interface{}) {
	stream := os.Stdout
	_, _ = whiteColor.Add(color.Bold).Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
}

// Infof format and print values to os.Stdout in white color
func Infof(format string, value ...interface{}) {
	Info(fmt.Sprintf(format, value...))
}

// Warn print values to os.Stdout in yellow color
func Warn(value ...interface{}) {
	stream := os.Stdout
	_, _ = yellowColor.Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
}

// Warnf format and print values to os.Stdout in yellow color
func Warnf(format string, value ...interface{}) {
	Warn(fmt.Sprintf(format, value...))
}

// Error print values to os.Stderr in red color
func Error(value ...interface{}) {
	stream := os.Stderr
	_, _ = redColor.Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
}

// Errorf format and print values to os.Stderr in red color
func Errorf(format string, value ...interface{}) {
	Error(fmt.Sprintf(format, value...))
}

// Fatal print values to os.Stderr in red color
// and invoke os.Exit with non zero status code
func Fatal(value ...interface{}) {
	stream := os.Stderr
	_, _ = redColor.Add(color.Bold).Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
	os.Exit(1)
}

// Fatalf format and print values to os.Stderr in red color
// and invoke os.Exit with non zero status code
func Fatalf(format string, value ...interface{}) {
	Fatal(fmt.Sprintf(format, value...))
}

// DieIf invoke os.Exit if err is not nil
func DieIf(err error) {
	if err != nil {
		Fatalf("\nCommand failed: %v", err)
	}
}
