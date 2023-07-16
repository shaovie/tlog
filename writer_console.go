package tlog

import (
    "os"
    "sync"
)

type WriteToConsole struct {
    mtx sync.Mutex
}

func NewWriteToConsole() *WriteToConsole { return &WriteToConsole{} }

func (w *WriteToConsole) Write(e Encoder, p []byte) (n int, err error) {
    w.mtx.Lock()
    defer w.mtx.Unlock()
    return os.Stdout.Write(p)
}
