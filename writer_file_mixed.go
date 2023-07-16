package tlog

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"syscall"
)

type WriteToFileMixed struct {
	newFileYear   int
	newFileMonth  int
	newFileDay    int
	fd            int
	dir           string
	logFilePrefix string
	fileStoreMode FileStoreModeT

	mtx sync.Mutex
}

func NewWriteToFileMixed(opts ...Option) Writer {
	opt := setOptions(opts...)

	w := &WriteToFileMixed{
		fd:            -1,
		dir:           opt.logDir,
		logFilePrefix: opt.logFilePrefix,
		fileStoreMode: opt.fileStoreMode,
	}
	if w.dir != "" {
		if err := os.MkdirAll(w.dir, 0755); err != nil {
			panic(errors.New("newlog mkdir fail! " + err.Error()))
		}
	}
	if w.fileStoreMode == AppendOneFile {
		if err := w.openAppendFile(); err != nil {
			panic(errors.New("newlog open file fail! " + err.Error()))
		}
	}
	return w
}
func (w *WriteToFileMixed) Write(e Encoder, p []byte) (n int, err error) {
	now := e.Now()
	year, month, day := now.Date()

	w.mtx.Lock()
	defer w.mtx.Unlock()

	if w.fileStoreMode == DailySplit {
		if err = w.newFile(year, int(month), day); err != nil {
			return
		}
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
func (w *WriteToFileMixed) newFile(year, month, day int) error {
	if w.newFileYear != year || w.newFileMonth != month || w.newFileDay != day {
		w.close()
		if err := w.openSeparateFile(year, month, day); err != nil {
			return err
		}
	}
	return nil
}
func (w *WriteToFileMixed) openAppendFile() (err error) {
	var fname string
	if len(w.logFilePrefix) == 0 {
		fname = fmt.Sprintf("%s.log", "tlog")
	} else {
		fname = fmt.Sprintf("%s.log", w.logFilePrefix)
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
func (w *WriteToFileMixed) openSeparateFile(year, month, day int) (err error) {
	var fname string
	if len(w.logFilePrefix) == 0 {
		fname = fmt.Sprintf("%d-%02d-%02d.log", year, month, day)
	} else {
		fname = fmt.Sprintf("%s-%d-%02d-%02d.log", w.logFilePrefix, year, month, day)
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
func (w *WriteToFileMixed) close() {
	if w.fd != -1 {
		syscall.Close(w.fd)
		w.fd = -1
	}
}
