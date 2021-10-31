package bytesize

import (
	"errors"
	"fmt"

	"github.com/kucjac/cleango/cgerrors"
)

// Bytes is the multiple-byte unit.
type Bytes int64

// MarshalText implements encoding.TextMarshaler interface.
func (b Bytes) MarshalText() ([]byte, error) {
	return []byte(b.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler interface.
func (b *Bytes) UnmarshalText(data []byte) error {
	tb, err := ParseBytes(string(data))
	if err != nil {
		return err
	}
	*b = tb
	return nil
}

// Exabytes gets the float value as metric exabytes.
func (b Bytes) Exabytes() float64 {
	v := b / Exabyte
	r := b % Exabyte
	return float64(v) + float64(r)/float64(Exabyte)
}

// Petabytes gets the float value as metric petabytes.
func (b Bytes) Petabytes() float64 {
	v := b / Petabyte
	r := b % Petabyte
	return float64(v) + float64(r)/float64(Petabyte)
}

// Terabytes gets the float value as metric terabytes.
func (b Bytes) Terabytes() float64 {
	v := b / Terabyte
	r := b % Terabyte
	return float64(v) + float64(r)/float64(Terabyte)
}

// Gigabytes gets the float value as metric gigabytes.
func (b Bytes) Gigabytes() float64 {
	v := b / Gigabyte
	r := b % Gigabyte
	return float64(v) + float64(r)/float64(Gigabyte)
}

// Megabytes gets the float value as metric megabytes.
func (b Bytes) Megabytes() float64 {
	v := b / Megabyte
	r := b % Megabyte
	return float64(v) + float64(r)/float64(Megabyte)
}

// Kilobytes gets the float value as metric kilobytes.
func (b Bytes) Kilobytes() float64 {
	v := b / Kilobyte
	r := b % Kilobyte
	return float64(v) + float64(r)/float64(Kilobyte)
}

func (b Bytes) String() string {
	switch {
	case b%Exabyte == 0:
		return fmt.Sprintf("%dEB", b/Exabyte)
	case b%Petabyte == 0:
		return fmt.Sprintf("%dPB", b/Petabyte)
	case b%Terabyte == 0:
		return fmt.Sprintf("%dTB", b/Terabyte)
	case b%Gigabyte == 0:
		return fmt.Sprintf("%dGB", b/Gigabyte)
	case b%Megabyte == 0:
		return fmt.Sprintf("%dMB", b/Megabyte)
	case b%Kilobyte == 0:
		return fmt.Sprintf("%dKB", b/Kilobyte)
	default:
		return fmt.Sprintf("%dB", b)
	}
}

// HumanReadable gets the human-readable
func (b Bytes) HumanReadable() string {
	switch {
	case b > Exabyte:
		return fmt.Sprintf("%.1f EB", b.Exabytes())
	case b > Petabyte:
		return fmt.Sprintf("%.1f PB", b.Petabytes())
	case b > Terabyte:
		return fmt.Sprintf("%.1f TB", b.Terabytes())
	case b > Gigabyte:
		return fmt.Sprintf("%.1f GB", b.Gigabytes())
	case b > Megabyte:
		return fmt.Sprintf("%.1f MB", b.Megabytes())
	case b > Kilobyte:
		return fmt.Sprintf("%.1f KB", b.Kilobytes())
	default:
		return fmt.Sprintf("%d B", b)
	}
}

// Byte is the lowest byte unit.
const Byte Bytes = 1

// Metric system representation
const (
	Kilobyte = 1000 * Byte
	Megabyte = 1000 * Kilobyte
	Gigabyte = 1000 * Megabyte
	Terabyte = 1000 * Gigabyte
	Petabyte = 1000 * Terabyte
	Exabyte  = 1000 * Petabyte
)

// IEC system representation.
const (
	Kibibyte = 1024 * Byte
	Mebibyte = 1024 * Kibibyte
	Gibibyte = 1024 * Mebibyte
	Tebibyte = 1024 * Gibibyte
	Pebibyte = 1024 * Tebibyte
	Exbibyte = 1024 * Pebibyte
)

// ParseBytes parses a unit string.
// Valid time units are "B", "kB", "KB", "MB", "GB", "TB", "PB", "EB", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"
// nolint: gocognit
func ParseBytes(s string) (Bytes, error) {
	// [-+]?([0-9]*(\.[0-9]*)?[a-z]+)+
	orig := s
	var d int64
	neg := false

	// Consume [-+]?
	if s != "" {
		c := s[0]
		if c == '-' || c == '+' {
			neg = c == '-'
			s = s[1:]
		}
	}
	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return 0, nil
	}
	if s == "" {
		return 0, cgerrors.ErrInvalidArgument("byte size: invalid unit " + quote(orig))
	}
	for s != "" {
		var (
			v, f  int64       // integers before, after decimal point
			scale float64 = 1 // value = v + f/scale
		)

		var err error

		// The next character must be [0-9.]
		if !(s[0] == '.' || '0' <= s[0] && s[0] <= '9') {
			return 0, cgerrors.ErrInvalidArgument("byte size: invalid unit " + quote(orig))
		}
		// Consume [0-9]*
		pl := len(s)
		v, s, err = leadingInt(s)
		if err != nil {
			return 0, cgerrors.ErrInvalidArgument("byte size: invalid unit " + quote(orig))
		}
		pre := pl != len(s) // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		post := false
		if s != "" && s[0] == '.' {
			s = s[1:]
			pl := len(s)
			f, scale, s = leadingFraction(s)
			post = pl != len(s)
		}
		if !pre && !post {
			// no digits (e.g. ".s" or "-.s")
			return 0, cgerrors.ErrInvalidArgument("byte size: invalid unit " + quote(orig))
		}

		// Consume unit.
		unit := int64(Byte)
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if c == '.' || '0' <= c && c <= '9' {
				break
			}
		}
		if i != 0 {
			u := s[:i]
			s = s[i:]
			var ok bool
			unit, ok = unitMap[u]
			if !ok {
				return 0, cgerrors.ErrInvalidArgument("byte size: unknown unit " + quote(u) + " in byte size" + quote(orig))
			}
		}
		if v > (1<<63-1)/unit {
			// overflow
			return 0, cgerrors.ErrInvalidArgument("byte size: invalid unit " + quote(orig))
		}
		v *= unit
		if f > 0 {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
			v += int64(float64(f) * (float64(unit) / scale))
			if v < 0 {
				// overflow
				return 0, cgerrors.ErrInvalidArgument("byte size: invalid unit " + quote(orig))
			}
		}
		d += v
		if d < 0 {
			// overflow
			return 0, errors.New("time: invalid unit " + quote(orig))
		}
	}

	if neg {
		d = -d
	}
	return Bytes(d), nil
}

func quote(s string) string {
	buf := make([]byte, 1, len(s)+2) // slice will be at least len(s) + quotes
	buf[0] = '"'
	for i, c := range s {
		if c >= runeSelf || c < ' ' {
			// This means you are asking us to parse a time.Bytes or
			// time.Location with unprintable or non-ASCII characters in it.
			// We don't expect to hit this case very often. We could try to
			// reproduce strconv.Quote's behavior with full fidelity but
			// given how rarely we expect to hit these edge cases, speed and
			// conciseness are better.
			var width int
			if c == runeError {
				width = 1
				if i+2 < len(s) && s[i:i+3] == string(runeError) {
					width = 3
				}
			} else {
				width = len(string(c))
			}
			for j := 0; j < width; j++ {
				buf = append(buf, `\x`...)
				buf = append(buf, lowerhex[s[i+j]>>4])
				buf = append(buf, lowerhex[s[i+j]&0xF])
			}
		} else {
			if c == '"' || c == '\\' {
				buf = append(buf, '\\')
			}
			buf = append(buf, string(c)...)
		}
	}
	buf = append(buf, '"')
	return string(buf)
}

// These are borrowed from unicode/utf8 and strconv and replicate behavior in
// that package, since we can't take a dependency on either.
const (
	lowerhex  = "0123456789abcdef"
	runeSelf  = 0x80
	runeError = '\uFFFD'
)

var errLeadingInt = errors.New("time: bad [0-9]*") // never printed

// leadingInt consumes the leading [0-9]* from s.
func leadingInt(s string) (x int64, rem string, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > (1<<63-1)/10 {
			// overflow
			return 0, "", errLeadingInt
		}
		x = x*10 + int64(c) - '0'
		if x < 0 {
			// overflow
			return 0, "", errLeadingInt
		}
	}
	return x, s[i:], nil
}

// leadingFraction consumes the leading [0-9]* from s.
// It is used only for fractions, so does not return an error on overflow,
// it just stops accumulating precision.
func leadingFraction(s string) (x int64, scale float64, rem string) {
	i := 0
	scale = 1
	overflow := false
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if overflow {
			continue
		}
		if x > (1<<63-1)/10 {
			// It's possible for overflow to give a positive number, so take care.
			overflow = true
			continue
		}
		y := x*10 + int64(c) - '0'
		if y < 0 {
			overflow = true
			continue
		}
		x = y
		scale *= 10
	}
	return x, scale, s[i:]
}

var unitMap = map[string]int64{
	"B":   int64(Byte),
	"kB":  int64(Kilobyte),
	"KB":  int64(Kilobyte),
	"MB":  int64(Megabyte),
	"GB":  int64(Gigabyte), // U+00B5 = micro symbol
	"TB":  int64(Terabyte), // U+03BC = Greek letter mu
	"PB":  int64(Petabyte),
	"EB":  int64(Exabyte),
	"KiB": int64(Kibibyte),
	"MiB": int64(Mebibyte),
	"GiB": int64(Gibibyte),
	"TiB": int64(Tebibyte),
	"PiB": int64(Pebibyte),
	"EiB": int64(Exbibyte),
}
