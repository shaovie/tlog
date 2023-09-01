package tlog

import (
	"encoding/json"
)

type Options struct {
	// for json/text
	omitEmpty bool

	format int

	timeFormat int

	level int

	writer Writer

	// for output file
	logDir        string
	logFilePrefix string
	fileStoreMode FileStoreModeT

	anyMarshalFunc AnyMarshalFuncT

	// for simple post
	postUrl string
}

type Option func(*Options)

func setOptions(optL ...Option) *Options {
	opts := &Options{
		omitEmpty:      true,
		format:         FormatJson,
		level:          AllLevel,
		writer:         NewWriteToConsole(),
		timeFormat:     HumanReadableTimeMs,
		logDir:         "logs",
		logFilePrefix:  "tlog",
		fileStoreMode:  DailySplit,
		anyMarshalFunc: json.Marshal,
	}

	for _, opt := range optL {
		opt(opts)
	}
	return opts
}

// If you don't want to output anything, you can use io.Discard
func SetWriter(w Writer) Option {
	return func(o *Options) {
		if w == nil {
			panic("tlog:SetWriter param is illegal")
		}
		o.writer = w
	}
}

// Set prefix `time` format
func TimeFormat(v int) Option {
	if v < 1 {
		panic("tlog:TimeFormat param is illegal")
	}
	return func(o *Options) {
		o.timeFormat = v
	}
}

// json/text
func Format(v int) Option {
	if v != FormatJson && v != FormatText {
		panic("tlog:Format param is illegal")
	}
	return func(o *Options) {
		o.format = v
	}
}

// for json:string,array
func OmitEmpty(v bool) Option {
	return func(o *Options) {
		o.omitEmpty = v
	}
}

// for output file
func LogDir(v string) Option {
	if len(v) == 0 {
		panic("tlog:LogDir param is illegal")
	}
	return func(o *Options) {
		o.logDir = v
	}
}
func LogFilePrefix(v string) Option {
	if len(v) == 0 {
		panic("tlog:LogFilePrefix param is illegal")
	}
	return func(o *Options) {
		o.logFilePrefix = v
	}
}
func FileStoreMode(v FileStoreModeT) Option {
	if v < 1 {
		panic("tlog:FileStoreMode param is illegal")
	}
	return func(o *Options) {
		o.fileStoreMode = v
	}
}

// for simple post
func PostUrl(v string) Option {
	if len(v) == 0 {
		panic("tlog:PostUrl param is illegal")
	}
	return func(o *Options) {
		o.postUrl = v
	}
}

func AnyMarshalFunc(f AnyMarshalFuncT) Option {
	if f == nil {
		panic("tlog:AnyMarshalFuncT param is illegal")
	}
	return func(o *Options) {
		o.anyMarshalFunc = f
	}
}
