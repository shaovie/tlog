package tlog

import (
    "fmt"
	"strconv"
    "reflect"
    "unsafe"
	"time"
)

type encoderText struct {
	encoder

    omitEmpty bool
	timeFormat int
}

func (e *encoderText) init() {
	e.now = time.Now()
	e.buf = e.buf[:0]
	e.appendHeaderTime()
    switch e.level {
    case DebugLevel:
        e.buf = append(e.buf, " debug"...)
    case InfoLevel:
        e.buf = append(e.buf, " info"...)
    case WarnLevel:
        e.buf = append(e.buf, " warn"...)
    case ErrorLevel:
        e.buf = append(e.buf, " error"...)
    case FatalLevel:
        e.buf = append(e.buf, " fatal"...)
    case PanicLevel:
        e.buf = append(e.buf, " panic"...)
    }
}
func (e *encoderText) OmitEmpty(v bool) Encoder { 
    e.omitEmpty = v
    return e
}
func (e *encoderText) appendKey(k string) {
	e.buf = append(e.buf, ' ')
	e.appendString(k)
    e.buf = append(e.buf, '=')
}
func (e *encoderText) fastAppendKey(k string) {
	e.buf = append(e.buf, ' ')
	e.fastAppendString(k)
    e.buf = append(e.buf, '=')
}
func (e *encoderText) appendHeaderTime() {
	if e.timeFormat == HumanReadableTime {
		e.appendHumanReadableTime()
	} else if e.timeFormat == HumanReadableTimeMs {
		e.appendHumanReadableTimeMs()
	} else if e.timeFormat == UnixTimestamp {
		e.buf = strconv.AppendInt(e.buf, e.now.Unix(), 10)
	} else if e.timeFormat == UnixTimestamp {
		e.buf = strconv.AppendInt(e.buf, e.now.Unix(), 10)
	} else if e.timeFormat == UnixTimestampMs {
		e.buf = strconv.AppendInt(e.buf, e.now.UnixMilli(), 10)
	}
}
func (e *encoderText) Fmt(k, format string, v ...any) Encoder {
	if e == nil {
		return nil
	}
    bf := make([]byte, 0, 256) // TODO
	bf = fmt.Appendf(bf, format, v...)
    if e.omitEmpty && len(bf) == 0 {
        return e
    }
	e.appendKey(k)
    if len(bf) > 0 {
        e.appendString(*(*string)(unsafe.Pointer(&bf)))
    }
    return e
}
func (e *encoderText) FastStr(k, v string) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(v) == 0 {
        return e
    }
	e.fastAppendKey(k)
    e.fastAppendString(v)
    return e
}
func (e *encoderText) Str(k, v string) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(v) == 0 {
        return e
    }
	e.appendKey(k)
	e.appendString(v)
    return e
}
func (e *encoderText) Strs(k string, vals []string) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(vals) == 0 {
        return e
    }
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
	}
    e.buf = append(e.buf, '[')
	e.buf = append(e.buf, '"')
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
    return e
}
func (e *encoderText) Bool(k string, v bool) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendBool(e.buf, v)
	return e
}
func (e *encoderText) Bools(k string, vals []bool) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(vals) == 0 {
        return e
    }
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
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
    return e
}
func (e *encoderText) Int(k string, v int) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendInt(e.buf, int64(v), 10)
	return e
}
func (e *encoderText) Ints(k string, vals []int) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(vals) == 0 {
        return e
    }
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
	}
    e.buf = append(e.buf, '[')
    e.appendInts(vals)
    e.buf = append(e.buf, ']')
    return e
}
func (e *encoderText) Int8(k string, v int8) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendInt(e.buf, int64(v), 10)
	return e
}
func (e *encoderText) Ints8(k string, vals []int8) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(vals) == 0 {
        return e
    }
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
	}
    e.buf = append(e.buf, '[')
    e.appendInts8(vals)
    e.buf = append(e.buf, ']')
    return e
}
func (e *encoderText) Int16(k string, v int16) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendInt(e.buf, int64(v), 10)
	return e
}
func (e *encoderText) Ints16(k string, vals []int16) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(vals) == 0 {
        return e
    }
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
	}
    e.buf = append(e.buf, '[')
    e.appendInts16(vals)
    e.buf = append(e.buf, ']')
    return e
}
func (e *encoderText) Int32(k string, v int32) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendInt(e.buf, int64(v), 10)
	return e
}
func (e *encoderText) Ints32(k string, vals []int32) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(vals) == 0 {
        return e
    }
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
	}
    e.buf = append(e.buf, '[')
    e.appendInts32(vals)
    e.buf = append(e.buf, ']')
    return e
}
func (e *encoderText) Int64(k string, v int64) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendInt(e.buf, v, 10)
	return e
}
func (e *encoderText) Ints64(k string, vals []int64) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
	}
    e.buf = append(e.buf, '[')
    e.appendInts64(vals)
    e.buf = append(e.buf, ']')
    return e
}
func (e *encoderText) Uint(k string, v uint) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendUint(e.buf, uint64(v), 10)
	return e
}
func (e *encoderText) Uints(k string, vals []uint) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(vals) == 0 {
        return e
    }
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
	}
    e.buf = append(e.buf, '[')
    e.appendUints(vals)
    e.buf = append(e.buf, ']')
    return e
}
func (e *encoderText) Uint8(k string, v uint8) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendUint(e.buf, uint64(v), 10)
	return e
}
func (e *encoderText) Uints8(k string, vals []uint8) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(vals) == 0 {
        return e
    }
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
	}
    e.buf = append(e.buf, '[')
    e.appendUints8(vals)
    e.buf = append(e.buf, ']')
    return e
}
func (e *encoderText) Uint16(k string, v uint16) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendUint(e.buf, uint64(v), 10)
	return e
}
func (e *encoderText) Uints16(k string, vals []uint16) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(vals) == 0 {
        return e
    }
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
	}
    e.buf = append(e.buf, '[')
    e.appendUints16(vals)
    e.buf = append(e.buf, ']')
    return e
}
func (e *encoderText) Uint32(k string, v uint32) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendUint(e.buf, uint64(v), 10)
	return e
}
func (e *encoderText) Uints32(k string, vals []uint32) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(vals) == 0 {
        return e
    }
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
	}
    e.buf = append(e.buf, '[')
    e.appendUints32(vals)
    e.buf = append(e.buf, ']')
    return e
}
func (e *encoderText) Uint64(k string, v uint64) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = strconv.AppendUint(e.buf, v, 10)
	return e
}
func (e *encoderText) Uints64(k string, vals []uint64) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(vals) == 0 {
        return e
    }
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
	}
    e.buf = append(e.buf, '[')
    e.appendUints64(vals)
    e.buf = append(e.buf, ']')
    return e
}
func (e *encoderText) Float32(k string, v float32) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
    e.appendFloat(float64(v), 32)
	return e
}
func (e *encoderText) Floats32(k string, vals []float32) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(vals) == 0 {
        return e
    }
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
	}
    e.buf = append(e.buf, '[')
    e.appendFloats32(vals)
    e.buf = append(e.buf, ']')
    return e
}
func (e *encoderText) Float64(k string, v float64) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
    e.appendFloat(v, 64)
	return e
}
func (e *encoderText) Floats64(k string, vals []float64) Encoder {
	if e == nil {
		return nil
	}
    if e.omitEmpty && len(vals) == 0 {
        return e
    }
	e.appendKey(k)
	if vals == nil {
		e.buf = append(e.buf, 'n', 'u', 'l', 'l')
        return e
    }
	if len(vals) == 0 {
		e.buf = append(e.buf, '[', ']')
        return e
	}
    e.buf = append(e.buf, '[')
    e.appendFloats64(vals)
    e.buf = append(e.buf, ']')
    return e
}
func (e *encoderText) Type(k string, v any) Encoder {
	if e == nil {
		return nil
	}
    if v == nil {
        return e.Str(k, "<nil>")
    }
	return e.Str(k, reflect.TypeOf(v).String())
}
func (e *encoderText) Time(k string, t time.Time, format string) Encoder {
	if e == nil {
		return nil
	}
	e.appendKey(k)
	e.buf = append(e.buf, '"')
	e.appendTime(t, format)
	e.buf = append(e.buf, '"')
	return e
}
func (e *encoderText) RawJSON(k string, b []byte) Encoder {
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
func (e *encoderText) Msg(s string) {
	if e == nil {
		return
	}
    if len(s) > 0 {
        e.Str("msg", s)
    }
    e.write()
}
func (e *encoderText) Msgf(format string, v ...any) {
	if e == nil {
		return
	}
    e.Fmt("msg", format, v...)
    e.write()
}
func (e *encoderText) Go() {
	if e == nil {
		return
	}
    e.write()
}
func (e *encoderText) write() {
	e.buf = append(e.buf, '\n')
	e.writer.Write(e, e.buf)
}
