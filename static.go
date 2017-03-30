// Copyright (c) 2017 Henry Slawniak <https://henry.computer/>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func serveFile(w http.ResponseWriter, r *http.Request, path string) {
	if path == "./client/" {
		path = "./client/index.html"
	}
	content, sum, mod, err := readFile(path)
	if err != nil {
		http.Error(w, "Could not read file", http.StatusInternalServerError)
		fmt.Printf("%s:%s\n", path, err.Error())
		return
	}
	mime := mime.TypeByExtension(filepath.Ext(path))
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Cache-Control", "public, no-cache")
	w.Header().Set("Last-Modified", mod.Format(time.RFC1123))
	if r.Header.Get("If-None-Match") == sum {
		w.WriteHeader(http.StatusNotModified)
		w.Header().Set("ETag", sum)
		return
	}
	w.Header().Set("ETag", sum)
	w.Write(content)
}

func readFile(path string) ([]byte, string, time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", time.Now(), err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, "", time.Now(), err
	}

	cont, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, "", time.Now(), err
	}

	return cont, fmt.Sprintf("%x", md5.Sum(cont)), stat.ModTime(), nil
}
