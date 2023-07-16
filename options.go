package tlog

import (
)

type Options struct {
    // for json/text
    omitEmpty bool

	format int

	timeFormat int

	level int

	writer Writer

    // for output file
    logDir string
    logFilePrefix string
    fileStoreMode FileStoreModeT

    // for simple post
    postUrl string
}

type Option func(*Options)

func setOptions(optL ...Option) *Options {
	opts := &Options{
        omitEmpty:  true,
		format:     FormatJson,
		level:      AllLevel,
		writer:     NewWriteToConsole(),
		timeFormat: HumanReadableTimeMs,
        logDir: "logs",
        logFilePrefix: "tlog",
        fileStoreMode: DailySplit,
	}

	for _, opt := range optL {
		opt(opts)
	}
	return opts
}

// If you don't want to output anything, you can use io.Discard
func SetWriter(w Writer) Option {
	return func(o *Options) {
		if w != nil {
			o.writer = w
		}
	}
}
// Set prefix `time` format
func TimeFormat(v int) Option {
	return func(o *Options) {
		if v > 0 {
			o.timeFormat = v
		}
	}
}
// json/text
func Format(v int) Option {
	return func(o *Options) {
		if v == FormatJson || v == FormatText {
			o.format = v
		}
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
	return func(o *Options) {
        if len(v) > 0 {
            o.logDir = v
        }
	}
}
func LogFilePrefix(v string) Option {
	return func(o *Options) {
        if len(v) > 0 {
            o.logFilePrefix = v
        }
	}
}
func FileStoreMode(v FileStoreModeT) Option {
	return func(o *Options) {
        if v > 0 {
            o.fileStoreMode = v
        }
	}
}

// for simple post
func PostUrl(v string) Option {
	return func(o *Options) {
        if len(v) > 0 {
            o.postUrl = v
        }
	}
}
