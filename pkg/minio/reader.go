package minio

import (
	"fmt"
	"io"
)

type ProgressReader struct {
	reader  io.Reader
	total   int64
	written int64
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	pr.written += int64(n)
	progress := float64(pr.written) / float64(pr.total) * 100
	fmt.Println("progress:", progress) // 这里打印进度
	return
}
