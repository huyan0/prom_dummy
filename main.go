package main

import (
	"compress/gzip"
	"io"
	"net/http"
	"sync"
)

const (
	contentTypeHeader     = "Content-Type"
	contentEncodingHeader = "Content-Encoding"
	acceptEncodingHeader  = "Accept-Encoding"
	contentType           = "text/plain; version=0.0.4; charset=utf-8"
)

var gzipPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(nil)
	},
}

func main() {
	text := `# HELP a_counter Counts things
# TYPE a_counter counter
a_counter{R="V",key="value"} 33
a_counter{R="V",key1="value1"} 33
a_counter{R="V",key="value1"} 33
# HELP a_valuerecorder Records values
# TYPE a_valuerecorder histogram
a_valuerecorder_bucket{R="V",key="value",le="+Inf"} 8
a_valuerecorder_sum{R="V",key="value"} 2900
a_valuerecorder_count{R="V",key="value"} 8`
	dH := &dummyHanlder{text: text}

	http.Handle("/", dH)
	http.ListenAndServe(":8888", nil)

}

type dummyHanlder struct {
	text string
}

// code from promhttp
func (d *dummyHanlder) ServeHTTP(rsp http.ResponseWriter, r *http.Request) {
	header := rsp.Header()
	// type is text
	header.Set(contentTypeHeader, string(contentType))
	w := io.Writer(rsp)

	// gzip compression
	header.Set(contentEncodingHeader, "gzip")
	gz := gzipPool.Get().(*gzip.Writer)
	defer gzipPool.Put(gz)

	gz.Reset(w)
	defer gz.Close()

	w = gz

	w.Write([]byte(d.text))
}
