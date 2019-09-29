// Copyright (c) 2019 Fadhli Dzil Ikram. All rights reserved.
// Use of this source code is governed by a MIT license that can be found in
// the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

type lineWriter struct {
	dst io.Writer

	mu  sync.Mutex
	buf bytes.Buffer
}

func (w *lineWriter) Write(p []byte) (int, error) {
	var offset int
	var err error
	w.mu.Lock()
	for err == nil && offset < len(p) {
		var noNewLine bool
		cp := p[offset:]
		if pos := bytes.IndexByte(cp, '\n'); pos >= 0 {
			cp = cp[:pos+1]
		} else {
			noNewLine = true
		}
		w.buf.Write(cp)
		offset += len(cp)
		if noNewLine {
			break
		}
		_, err = w.dst.Write(w.buf.Bytes())
		w.buf.Reset()
	}
	w.mu.Unlock()
	return offset, err
}

// LineWriter returns new writer that buffer any writes until new line and
// write them to dst.
func LineWriter(dst io.Writer) io.Writer {
	return &lineWriter{dst: dst}
}

type writerFunc func([]byte) (int, error)

func (fn writerFunc) Write(p []byte) (int, error) {
	return fn(p)
}

// PrefixWriter returns new writer that adds prefix to every call to io.Write.
func PrefixWriter(dst io.Writer, prefix string) io.Writer {
	return writerFunc(func(p []byte) (int, error) {
		return fmt.Fprintf(dst, "%s%s", prefix, p)
	})
}
