package tlog

import (
	"os"
	"sync"
)

const (
	DebugLevel int = 1 << 0
	InfoLevel  int = 1 << 1
	WarnLevel  int = 1 << 2
	ErrorLevel int = 1 << 3
	FatalLevel int = 1 << 4
	PanicLevel int = 1 << 5
	AllLevel   int = (DebugLevel | InfoLevel | WarnLevel | ErrorLevel | FatalLevel | PanicLevel)

	FormatJson int = 1
	FormatText int = 2
)

type TLog struct {
	omitEmpty      bool // for json
	format         int
	level          int
	timeFormat     int
	anyMarshalFunc AnyMarshalFuncT

	encoderTextPool sync.Pool
	encoderJsonPool sync.Pool
	writer          Writer
}

func New(opts ...Option) *TLog {
	opt := setOptions(opts...)

	tl := &TLog{
		omitEmpty:      opt.omitEmpty,
		format:         opt.format,
		level:          opt.level,
		writer:         opt.writer,
		timeFormat:     opt.timeFormat,
		anyMarshalFunc: opt.anyMarshalFunc,
		encoderTextPool: sync.Pool{
			New: func() any {
				return &encoderText{
					encoder: encoder{
						buf: make([]byte, 0, 512),
					},
				}
			},
		},
		encoderJsonPool: sync.Pool{
			New: func() any {
				return &encoderJson{
					encoder: encoder{
						buf: make([]byte, 0, 512),
					},
				}
			},
		},
	}

	return tl
}

func (tl *TLog) newEncoder(lvl int, doneCallback func(s string)) Encoder {
	if tl.level&lvl == 0 {
		if doneCallback != nil {
			doneCallback("(level diabled)")
		}
		return nil
	}
	var e Encoder
	if tl.format == FormatJson {
		obj := tl.encoderJsonPool.Get().(*encoderJson)
		obj.level = lvl
		obj.omitEmpty = tl.omitEmpty
		obj.timeFormat = tl.timeFormat
		obj.writer = tl.writer
		obj.doneCallback = doneCallback
		obj.anyMarshalFunc = tl.anyMarshalFunc
		obj.init()
		e = obj
	} else if tl.format == FormatText {
		obj := tl.encoderTextPool.Get().(*encoderText)
		obj.level = lvl
		obj.omitEmpty = tl.omitEmpty
		obj.timeFormat = tl.timeFormat
		obj.writer = tl.writer
		obj.doneCallback = doneCallback
		obj.anyMarshalFunc = tl.anyMarshalFunc
		obj.init()
		e = obj
	}
	return e
}
func (tl *TLog) Debug() Encoder {
	return tl.newEncoder(DebugLevel, nil)
}
func (tl *TLog) Info() Encoder {
	return tl.newEncoder(InfoLevel, nil)
}
func (tl *TLog) Warn() Encoder {
	return tl.newEncoder(WarnLevel, nil)
}
func (tl *TLog) Error() Encoder {
	return tl.newEncoder(ErrorLevel, nil)
}
func (tl *TLog) Fatal() Encoder {
	return tl.newEncoder(FatalLevel, func(msg string) { os.Exit(1) })
}
func (tl *TLog) Panic() Encoder {
	return tl.newEncoder(PanicLevel, func(msg string) { panic(msg) })
}
