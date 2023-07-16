package tlog

import (
	"bytes"
	"net/http"
	"sync"
	"time"
)

type WriteToSimplePost struct {
	mtx sync.Mutex
	url string
}

func NewWriteToSimplePost(opts ...Option) *WriteToSimplePost {
	opt := setOptions(opts...)
	if len(opt.postUrl) == 0 {
		panic("newlog simple post, post url is empty")
	}
	return &WriteToSimplePost{url: opt.postUrl}
}

func (w *WriteToSimplePost) Write(e Encoder, p []byte) (n int, err error) {
	client := &http.Client{Timeout: 3000 * time.Millisecond}
	buffer := bytes.NewBuffer(p)
	request, err := http.NewRequest("POST", w.url, buffer)
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil && resp == nil {
		return 0, err
	}
	if resp != nil {
		defer resp.Body.Close()
	}
	return len(p), nil
}
