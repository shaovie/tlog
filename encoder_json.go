package tlog

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

type encoderJson struct {
	encoder
}

func (e *encoderJson) init() {
	e.now = time.Now()
	e.buf = e.buf[:0]
	e.buf = append(e.buf, '{')
	e.appendHeaderTime()
	switch e.level {
	case DebugLevel:
		e.FastStr("level", "debug")
	case InfoLevel:
		e.FastStr("level", "info")
	case WarnLevel:
		e.FastStr("level", "warn")
	case ErrorLevel:
		e.FastStr("level", "error")
	case FatalLevel:
		e.FastStr("level", "fatal")
	case PanicLevel:
		e.FastStr("level", "panic")
	}
}
func (e *encoderJson) OmitEmpty(v bool) Encoder {
	e.omitEmpty = v
	return e
}
func (e *encoderJson) AnyMarshalFunc(f AnyMarshalFuncT) Encoder {
	e.anyMarshalFunc = f
	return e
}
func (e *encoderJson) appendKey(k string) {
	if e.buf[len(e.buf)-1] != '{' {
		e.buf = append(e.buf, ',')
	}
	e.buf = append(e.buf, '"')
	e.appendString(k)
	e.buf = append(e.buf, '"', ':')
}
func (e *encoderJson) fastAppendKey(k string) {
	if e.buf[len(e.buf)-1] != '{' {
		e.buf = append(e.buf, ',')
	}
	e.buf = append(e.buf, '"')
	e.fastAppendString(k)
	e.buf = append(e.buf, '"', ':')
}
func (e *encoderJson) appendHeaderTime() {
	e.appendKey("time")
	if e.timeFormat == HumanReadableTime {
		e.buf = append(e.buf, '"')
		e.appendHumanReadableTime()
		e.buf = append(e.buf, '"')
	} else if e.timeFormat == HumanReadableTimeMs {
		e.buf = append(e.buf, '"')
		e.appendHumanReadableTimeMs()
		e.buf = append(e.buf, '"')
	} else if e.timeFormat == UnixTimestamp {
		e.buf = strconv.AppendInt(e.buf, e.now.Unix(), 10)
	} else if e.timeFormat == UnixTimestampMs {
		e.buf = strconv.AppendInt(e.buf, e.now.UnixMilli(), 10)
	}
}
func (e *encoderJson) Fmt(k, format string, v ...any) Encoder {
	if e == nil {
		return nil
	}
	bf := make([]byte, 0, 256) // TODO
	bf = fmt.Appendf(bf, format, v...)
	if e.omitEmpty && len(bf) == 0 {
		return e
	}
	e.appendKey(k)
	e.buf = append(e.buf, '"')
	if len(bf) > 0 {
		e.appendString(*(*string)(unsafe.Pointer(&bf)))
	}
	e.buf = append(e.buf, '"')
	return e
}
func (e *encoderJson) FastStr(k, v string) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(v) == 0 {
		return e
	}
	e.fastAppendKey(k)
	e.buf = append(e.buf, '"')
	e.fastAppendString(v)
	e.buf = append(e.buf, '"')
	return e
}
func (e *encoderJson) Str(k, v string) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(v) == 0 {
		return e
	}
	e.appendKey(k)
	e.buf = append(e.buf, '"')
	e.appendString(v)
	e.buf = append(e.buf, '"')
	return e
}
func (e *encoderJson) Strs(k string, vals []string) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(vals) == 0 {
		return e
	}
	e.appendKey(k)
	e.appendStrings(vals)
	return e
}
func (e *encoderJson) Bool(k string, v bool) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendBool(e.buf, v)
	return e
}
func (e *encoderJson) Bools(k string, vals []bool) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(vals) == 0 {
		return e
	}
	e.appendKey(k)
	e.appendBools(vals)
	return e
}
func (e *encoderJson) Int(k string, v int) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendInt(e.buf, int64(v), 10)
	return e
}
func (e *encoderJson) Ints(k string, vals []int) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(vals) == 0 {
		return e
	}
	e.appendKey(k)
	e.appendInts(vals)
	return e
}
func (e *encoderJson) Int8(k string, v int8) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendInt(e.buf, int64(v), 10)
	return e
}
func (e *encoderJson) Ints8(k string, vals []int8) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(vals) == 0 {
		return e
	}
	e.appendKey(k)
	e.appendInts8(vals)
	return e
}
func (e *encoderJson) Int16(k string, v int16) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendInt(e.buf, int64(v), 10)
	return e
}
func (e *encoderJson) Ints16(k string, vals []int16) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(vals) == 0 {
		return e
	}
	e.appendKey(k)
	e.appendInts16(vals)
	return e
}
func (e *encoderJson) Int32(k string, v int32) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendInt(e.buf, int64(v), 10)
	return e
}
func (e *encoderJson) Ints32(k string, vals []int32) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(vals) == 0 {
		return e
	}
	e.appendKey(k)
	e.appendInts32(vals)
	return e
}
func (e *encoderJson) Int64(k string, v int64) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendInt(e.buf, v, 10)
	return e
}
func (e *encoderJson) Ints64(k string, vals []int64) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.appendInts64(vals)
	return e
}
func (e *encoderJson) Uint(k string, v uint) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendUint(e.buf, uint64(v), 10)
	return e
}
func (e *encoderJson) Uints(k string, vals []uint) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(vals) == 0 {
		return e
	}
	e.appendKey(k)
	e.appendUints(vals)
	return e
}
func (e *encoderJson) Uint8(k string, v uint8) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendUint(e.buf, uint64(v), 10)
	return e
}
func (e *encoderJson) Uints8(k string, vals []uint8) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(vals) == 0 {
		return e
	}
	e.appendKey(k)
	e.appendUints8(vals)
	return e
}
func (e *encoderJson) Uint16(k string, v uint16) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendUint(e.buf, uint64(v), 10)
	return e
}
func (e *encoderJson) Uints16(k string, vals []uint16) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(vals) == 0 {
		return e
	}
	e.appendKey(k)
	e.appendUints16(vals)
	return e
}
func (e *encoderJson) Uint32(k string, v uint32) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendUint(e.buf, uint64(v), 10)
	return e
}
func (e *encoderJson) Uints32(k string, vals []uint32) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(vals) == 0 {
		return e
	}
	e.appendKey(k)
	e.appendUints32(vals)
	return e
}
func (e *encoderJson) Uint64(k string, v uint64) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendUint(e.buf, v, 10)
	return e
}
func (e *encoderJson) Uints64(k string, vals []uint64) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(vals) == 0 {
		return e
	}
	e.appendKey(k)
	e.appendUints64(vals)
	return e
}
func (e *encoderJson) Float32(k string, v float32) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.appendFloat(float64(v), 32)
	return e
}
func (e *encoderJson) Floats32(k string, vals []float32) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(vals) == 0 {
		return e
	}
	e.appendKey(k)
	e.appendFloats32(vals)
	return e
}
func (e *encoderJson) Float64(k string, v float64) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.appendFloat(v, 64)
	return e
}
func (e *encoderJson) Floats64(k string, vals []float64) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(vals) == 0 {
		return e
	}
	e.appendKey(k)
	e.appendFloats64(vals)
	return e
}
func (e *encoderJson) Type(k string, v any) Encoder {
	if e == nil {
		return nil
	}
	if v == nil {
		return e.Str(k, "<nil>")
	}
	return e.Str(k, reflect.TypeOf(v).String())
}
func (e *encoderJson) Any(k string, v any) Encoder {
	if e == nil {
		return nil
	}
	marshaled, err := e.anyMarshalFunc(v)
	if err != nil {
		return e.Str(k, fmt.Sprintf("marshaling error: %s", err.Error()))
	}
	e.appendKey(k)
	e.buf = append(e.buf, marshaled...)
	return e
}
func (e *encoderJson) Time(k string, t time.Time, format string) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = append(e.buf, '"')
	e.appendTime(t, format)
	e.buf = append(e.buf, '"')
	return e
}
func (e *encoderJson) RawJSON(k string, b []byte) Encoder {
	if e == nil {
		return nil
	}
	if e.omitEmpty && len(b) == 0 {
		return e
	}
	e.appendKey(k)
	if b == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
		return e
	}
	e.appendRawJSON(k, b)
	return e
}
func (e *encoderJson) Msg(s string) {
	if e == nil {
		return
	}
	if len(s) > 0 {
		e.Str("msg", s)
	}
    e.Go()
}
func (e *encoderJson) Msgf(format string, v ...any) {
	if e == nil {
		return
	}
	e.Fmt("msg", format, v...)
    e.Go()
}
func (e *encoderJson) Go() {
	if e == nil {
		return
	}
	e.buf = append(e.buf, '}')
	e.buf = append(e.buf, '\n')
	e.writer.Write(e, e.buf)

    // Proper usage of a sync.Pool requires each entry to have approximately
	// the same memory cost. To obtain this property when the stored type
	// contains a variably-sized buffer, we add a hard limit on the maximum buffer
	// to place back in the pool.
	//
	// See https://golang.org/issue/23199
	if cap(e.buf) > (1 << 14) { // 16KiB
		return
	}
    e.tl.encoderJsonPool.Put(e)
}
