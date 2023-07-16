package tlog

import (
	"math"
	"strconv"
	"time"
	"unicode/utf8"
)

const (
	HumanReadableTime   int = 1 // 2023-07-14 21:08:20
	HumanReadableTimeMs int = 2 // 2023-07-14 21:08:20.212
	UnixTimestamp       int = 3 // 1689340100
	UnixTimestampMs     int = 4 // 1689340100123 // millisecond
)

type Encoder interface {
	// For writer
	Level() int
	Now() time.Time

	//=
	Fmt(k, format string, v ...any) Encoder

	Str(k, v string) Encoder
	Strs(k string, v []string) Encoder
	FastStr(k, v string) Encoder

	Bool(k string, v bool) Encoder
	Bools(k string, v []bool) Encoder

	Int(k string, v int) Encoder
	Ints(k string, v []int) Encoder
	Int8(k string, v int8) Encoder
	Ints8(k string, v []int8) Encoder
	Int16(k string, v int16) Encoder
	Ints16(k string, v []int16) Encoder
	Int32(k string, v int32) Encoder
	Ints32(k string, v []int32) Encoder
	Int64(k string, v int64) Encoder
	Ints64(k string, v []int64) Encoder

	Uint(k string, v uint) Encoder
	Uints(k string, v []uint) Encoder
	Uint8(k string, v uint8) Encoder
	Uints8(k string, v []uint8) Encoder
	Uint16(k string, v uint16) Encoder
	Uints16(k string, v []uint16) Encoder
	Uint32(k string, v uint32) Encoder
	Uints32(k string, v []uint32) Encoder
	Uint64(k string, v uint64) Encoder
	Uints64(k string, v []uint64) Encoder

	Float32(k string, v float32) Encoder
	Floats32(k string, v []float32) Encoder
	Float64(k string, v float64) Encoder
	Floats64(k string, v []float64) Encoder

	Type(k string, v any) Encoder

	Time(k string, t time.Time, format string) Encoder

	// RawJSON adds already encoded JSON to the log line under key.
	//
	// No sanity check is performed on b; it must not contain carriage returns and
	// be valid JSON
	RawJSON(k string, b []byte) Encoder

	// config
	OmitEmpty(v bool) Encoder

	// end
	Msg(s string)
	Msgf(format string, v ...any)
	Go()
}

type encoder struct {
	level int
	buf   []byte

	now    time.Time
	writer Writer

	doneCallback func(s string)
}

const hex = "0123456789abcdef"
const twoDigits = "00010203040506070809" +
	"10111213141516171819" +
	"20212223242526272829" +
	"30313233343536373839" +
	"40414243444546474849" +
	"50515253545556575859" +
	"60616263646566676869" +
	"70717273747576777879" +
	"80818283848586878889" +
	"90919293949596979899"

var noEscapeTable = [256]bool{}

func init() {
	for i := 0; i <= 0x7e; i++ {
		noEscapeTable[i] = i >= 0x20 && i != '\\' && i != '"'
	}
}

func (e *encoder) Level() int {
	return e.level
}
func (e *encoder) Now() time.Time {
	return e.now
}
func (e *encoder) appendTwoDigitsString(i int) {
	e.buf = append(e.buf, twoDigits[i*2:i*2+2]...)
}
func (e *encoder) appendHumanReadableTimeMs() {
	e.appendHumanReadableTime()
	e.buf = append(e.buf, '.')
	e.appendPaddingInt(e.now.Nanosecond()/1e6, 3)
}
func (e *encoder) appendHumanReadableTime() {
	year, month, day := e.now.Date()
	hour, min, sec := e.now.Clock()
	e.buf = append(e.buf, '2', '0')
	e.appendTwoDigitsString(year % 100)
	e.buf = append(e.buf, '-')
	e.appendTwoDigitsString(int(month))
	e.buf = append(e.buf, '-')
	e.appendTwoDigitsString(day)
	e.buf = append(e.buf, ' ')
	e.appendTwoDigitsString(hour)
	e.buf = append(e.buf, ':')
	e.appendTwoDigitsString(min)
	e.buf = append(e.buf, ':')
	e.appendTwoDigitsString(sec)
}
func (e *encoder) appendPaddingInt(i int, wid int) {
	// Assemble decimal in reverse order.
	var b [8]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	e.buf = append(e.buf, b[bp:]...)
}
func (e *encoder) fastAppendString(s string) {
	e.buf = append(e.buf, s...)
}
func (e *encoder) appendString(s string) {
	for i := 0; i < len(s); i++ {
		// Check if the character needs encoding. Control characters, slashes,
		// and the double quote need json encoding. Bytes above the ascii
		// boundary needs utf8 encoding.
		if !noEscapeTable[s[i]] {
			// We encountered a character that needs to be encoded. Switch
			// to complex version of the algorithm.
			e.appendStringComplex(s, i)
			return
		}
	}
	e.buf = append(e.buf, s...)
}

// appendStringComplex is used by appendString to take over an in
// progress JSON string encoding that encountered a character that needs
// to be encoded.
//
// Copied from github.com/rs/zerolog/internal/json/string.go
func (e *encoder) appendStringComplex(s string, i int) {
	start := 0
	for i < len(s) {
		b := s[i]
		if b >= utf8.RuneSelf {
			r, size := utf8.DecodeRuneInString(s[i:])
			if r == utf8.RuneError && size == 1 {
				// In case of error, first append previous simple characters to
				// the byte slice if any and append a replacement character code
				// in place of the invalid sequence.
				if start < i {
					e.buf = append(e.buf, s[start:i]...)
				}
				e.buf = append(e.buf, `\ufffd`...)
				i += size
				start = i
				continue
			}
			i += size
			continue
		}
		if noEscapeTable[b] {
			i++
			continue
		}
		// We encountered a character that needs to be encoded.
		// Let's append the previous simple characters to the byte slice
		// and switch our operation to read and encode the remainder
		// characters byte-by-byte.
		if start < i {
			e.buf = append(e.buf, s[start:i]...)
		}
		switch b {
		case '"', '\\':
			e.buf = append(e.buf, '\\', b)
		case '\b':
			e.buf = append(e.buf, '\\', 'b')
		case '\f':
			e.buf = append(e.buf, '\\', 'f')
		case '\n':
			e.buf = append(e.buf, '\\', 'n')
		case '\r':
			e.buf = append(e.buf, '\\', 'r')
		case '\t':
			e.buf = append(e.buf, '\\', 't')
		default:
			e.buf = append(e.buf, '\\', 'u', '0', '0', hex[b>>4], hex[b&0xF])
		}
		i++
		start = i
	}
	if start < len(s) {
		e.buf = append(e.buf, s[start:]...)
	}
}
func (e *encoder) appendStrings(vals []string) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[', '"')
	e.appendString(vals[0])
	e.buf = append(e.buf, '"')
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',', '"')
			e.appendString(val)
			e.buf = append(e.buf, '"')
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendBools(vals []bool) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[')
	e.buf = strconv.AppendBool(e.buf, vals[0])
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',')
			e.buf = strconv.AppendBool(e.buf, val)
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendInts(vals []int) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[')
	e.buf = strconv.AppendInt(e.buf, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',')
			e.buf = strconv.AppendInt(e.buf, int64(val), 10)
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendInts8(vals []int8) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[')
	e.buf = strconv.AppendInt(e.buf, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',')
			e.buf = strconv.AppendInt(e.buf, int64(val), 10)
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendInts16(vals []int16) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[')
	e.buf = strconv.AppendInt(e.buf, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',')
			e.buf = strconv.AppendInt(e.buf, int64(val), 10)
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendInts32(vals []int32) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[')
	e.buf = strconv.AppendInt(e.buf, int64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',')
			e.buf = strconv.AppendInt(e.buf, int64(val), 10)
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendInts64(vals []int64) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[')
	e.buf = strconv.AppendInt(e.buf, vals[0], 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',')
			e.buf = strconv.AppendInt(e.buf, val, 10)
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendUints(vals []uint) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[')
	e.buf = strconv.AppendUint(e.buf, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',')
			e.buf = strconv.AppendUint(e.buf, uint64(val), 10)
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendUints8(vals []uint8) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[')
	e.buf = strconv.AppendUint(e.buf, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',')
			e.buf = strconv.AppendUint(e.buf, uint64(val), 10)
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendUints16(vals []uint16) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[')
	e.buf = strconv.AppendUint(e.buf, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',')
			e.buf = strconv.AppendUint(e.buf, uint64(val), 10)
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendUints32(vals []uint32) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[')
	e.buf = strconv.AppendUint(e.buf, uint64(vals[0]), 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',')
			e.buf = strconv.AppendUint(e.buf, uint64(val), 10)
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendUints64(vals []uint64) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[')
	e.buf = strconv.AppendUint(e.buf, vals[0], 10)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',')
			e.buf = strconv.AppendUint(e.buf, val, 10)
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendFloats32(vals []float32) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[')
	e.appendFloat(float64(vals[0]), 32)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',')
			e.appendFloat(float64(val), 32)
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendFloats64(vals []float64) {
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return
	}
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
		return
	}
	e.buf = append(e.buf, '[')
	e.appendFloat(vals[0], 64)
	if len(vals) > 1 {
		for _, val := range vals[1:] {
			e.buf = append(e.buf, ',')
			e.appendFloat(val, 64)
		}
	}
	e.buf = append(e.buf, ']')
}
func (e *encoder) appendFloat(v float64, bitSize int) {
	if math.IsNaN(v) {
		e.buf = append(e.buf, `"NaN"`...)
	} else if math.IsInf(v, 1) {
		e.buf = append(e.buf, `"+Inf"`...)
	} else if math.IsInf(v, -1) {
		e.buf = append(e.buf, `"-Inf"`...)
	} else {
		e.buf = strconv.AppendFloat(e.buf, v, 'f', -1, bitSize)
	}
}
func (e *encoder) appendTime(t time.Time, format string) {
	e.buf = t.AppendFormat(e.buf, format)
}
func (e *encoder) appendRawJSON(k string, b []byte) {
	e.buf = append(e.buf, b...)
}
