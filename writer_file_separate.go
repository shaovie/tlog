package tlog

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"syscall"
)

type WriteToFileSeparate struct {
	debugWriter writeToFileSeparateLevel
	infoWriter  writeToFileSeparateLevel
	warnWriter  writeToFileSeparateLevel
	errorWriter writeToFileSeparateLevel
	fatalWriter writeToFileSeparateLevel
	panicWriter writeToFileSeparateLevel
}

func NewWriteToFileSeparate(opts ...Option) Writer {
	opt := setOptions(opts...)

	w := &WriteToFileSeparate{
		debugWriter: writeToFileSeparateLevel{
			fd:            -1,
			dir:           opt.logDir,
			logFilePrefix: opt.logFilePrefix,
			name:          "debug",
			fileStoreMode: opt.fileStoreMode,
		},
		infoWriter: writeToFileSeparateLevel{
			fd:            -1,
			dir:           opt.logDir,
			logFilePrefix: opt.logFilePrefix,
			name:          "info",
			fileStoreMode: opt.fileStoreMode,
		},
		warnWriter: writeToFileSeparateLevel{
			fd:            -1,
			dir:           opt.logDir,
			logFilePrefix: opt.logFilePrefix,
			name:          "warn",
			fileStoreMode: opt.fileStoreMode,
		},
		errorWriter: writeToFileSeparateLevel{
			fd:            -1,
			dir:           opt.logDir,
			logFilePrefix: opt.logFilePrefix,
			name:          "error",
			fileStoreMode: opt.fileStoreMode,
		},
		fatalWriter: writeToFileSeparateLevel{
			fd:            -1,
			dir:           opt.logDir,
			logFilePrefix: opt.logFilePrefix,
			name:          "fatal",
			fileStoreMode: opt.fileStoreMode,
		},
		panicWriter: writeToFileSeparateLevel{
			fd:            -1,
			dir:           opt.logDir,
			logFilePrefix: opt.logFilePrefix,
			name:          "panic",
			fileStoreMode: opt.fileStoreMode,
		},
	}
	if opt.logDir != "" {
		if err := os.MkdirAll(opt.logDir, 0755); err != nil {
			panic(errors.New("newlog mkdir fail! " + err.Error()))
		}
	}
	return w
}
func (w *WriteToFileSeparate) Write(e Encoder, p []byte) (n int, err error) {
	switch e.Level() {
	case DebugLevel:
		return w.debugWriter.Write(e, p)
	case InfoLevel:
		return w.infoWriter.Write(e, p)
	case WarnLevel:
		return w.warnWriter.Write(e, p)
	case ErrorLevel:
		return w.errorWriter.Write(e, p)
	case FatalLevel:
		return w.fatalWriter.Write(e, p)
	case PanicLevel:
		return w.panicWriter.Write(e, p)
	}
	return 0, nil
}

type writeToFileSeparateLevel struct {
	newFileYear   int
	newFileMonth  int
	newFileDay    int
	fd            int
	dir           string
	name          string
	logFilePrefix string
	fileStoreMode FileStoreModeT

	mtx sync.Mutex
}

func (w *writeToFileSeparateLevel) Write(e Encoder, p []byte) (n int, err error) {
	now := e.Now()
	year, month, day := now.Date()

	w.mtx.Lock()
	defer w.mtx.Unlock()

	if w.fileStoreMode == DailySplit {
		if err = w.newFile(year, int(month), day); err != nil {
			return
		}
	} else if w.fd == -1 {
		w.openAppendFile()
	}
	for {
		n, err = syscall.Write(w.fd, p)
		if err != nil && err == syscall.EINTR {
			continue
		}
		break
	}
	return
}
func (w *writeToFileSeparateLevel) newFile(year, month, day int) error {
	if w.newFileYear != year || w.newFileMonth != month || w.newFileDay != day {
		w.close()
		if err := w.openSeparateFile(year, month, day); err != nil {
			return err
		}
	}
	return nil
}
func (w *writeToFileSeparateLevel) openAppendFile() (err error) {
	var fname string
	if len(w.logFilePrefix) == 0 {
		fname = fmt.Sprintf("%s.log", w.name)
	} else {
		fname = fmt.Sprintf("%s-%s.log", w.logFilePrefix, w.name)
	}
	logFile := path.Join(w.dir, fname)
	for {
		w.fd, err = syscall.Open(logFile, syscall.O_CREAT|syscall.O_WRONLY|syscall.O_APPEND, 0644)
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			return err
		}
		break
	}
	return nil
}
func (w *writeToFileSeparateLevel) openSeparateFile(year, month, day int) (err error) {
	var fname string
	if len(w.logFilePrefix) == 0 {
		fname = fmt.Sprintf("%s-%d-%02d-%02d.log", w.name, year, month, day)
	} else {
		fname = fmt.Sprintf("%s-%s-%d-%02d-%02d.log", w.logFilePrefix, w.name, year, month, day)
	}
	logFile := path.Join(w.dir, fname)
	for {
		w.fd, err = syscall.Open(logFile, syscall.O_CREAT|syscall.O_WRONLY|syscall.O_APPEND, 0644)
		if err != nil {
			if err == syscall.EINTR {
				continue
			}
			return err
		}
		break
	}
	w.newFileYear, w.newFileMonth, w.newFileDay = year, month, day
	return nil
}
func (w *writeToFileSeparateLevel) close() {
	if w.fd != -1 {
		syscall.Close(w.fd)
		w.fd = -1
	}
}
