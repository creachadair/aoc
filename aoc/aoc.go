package aoc

import (
	"encoding"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

var doDebug = flag.Bool("debug", false, "Enable debug output")

// Dprintf acts as log.Printf if the --debug flag is set; otherwise it discards
// its input.
func Dprintf(msg string, args ...any) {
	if *doDebug {
		log.Printf(msg, args...)
	}
}

// MustReadInput reads the contents of the first command-line argument, or if
// none is specified it fully consumes the contents of stdin.  In case of
// error, MustReadInput calls log.Fatal.
func MustReadInput() []byte {
	var data []byte
	var err error
	if flag.NArg() == 0 {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(flag.Arg(0))
	}
	if err != nil {
		log.Fatalf("Read input: %v", err)
	}
	return data
}

// SplitLines splits input into lines.
func SplitLines(input []byte) []string {
	return strings.Split(strings.TrimSpace(string(input)), "\n")
}

// MustReadLines is shorthand for SplitLines(MustReadInput()).
func MustReadLines() []string { return SplitLines(MustReadInput()) }

// ParseInts parses strings as integers.
func ParseInts(ss []string) ([]int, error) {
	out := make([]int, len(ss))
	for i, s := range ss {
		v, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("value %d: %w", i+1, err)
		}
		out[i] = v
	}
	return out, nil
}

// Scanx matches the given re against input, and populates the specified
// arguments with the corresponding capturing groups. The number of capturing
// groups in re must match the number of args, and each arg must be a pointer
// to one of these types:
//
//	int      -- a (signed) decimal integer
//	[]int    -- multiple separated decimal integers (see below)
//	string   -- a single string
//	[]string -- a slice of separated strings (see below)
//	encoding.TextUnmarshaler
//
// Each match is decoded into the target argument, or an error is reported.
//
// By default, multiple values are separated by whitespace.
// Name the capture group "comma" to split on commas.
func Scanx(re *regexp.Regexp, input string, args ...any) error {
	if n := re.NumSubexp(); n != len(args) {
		return fmt.Errorf("want %d subexpressions, got %d", len(args), n)
	}
	m := re.FindStringSubmatch(input)
	if m == nil {
		return fmt.Errorf("input does not match %q", re)
	}
	for i, sub := range m[1:] {
		var err error
		switch arg := args[i].(type) {
		case *int:
			err = parseInt(sub, arg)
		case *[]int:
			*arg, err = ParseInts(splitFields(sub, re.SubexpNames()[i+1]))
		case *string:
			*arg = sub
		case *[]string:
			*arg = splitFields(sub, re.SubexpNames()[i+1])
		case encoding.TextUnmarshaler:
			err = arg.UnmarshalText([]byte(sub))
		default:
			err = fmt.Errorf("incompatible type %T", args[i])
		}
		if err != nil {
			return fmt.Errorf("argument %d: %w", i+1, err)
		}
	}
	return nil
}

func parseInt[T constraints.Integer](s string, vp *T) error {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	*vp = T(v)
	return nil
}

func splitFields(s, name string) []string {
	switch name {
	case "comma":
		return strings.Split(s, ",")
	default:
		return strings.Fields(s)
	}
}
