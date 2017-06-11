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
	"net/http"
	"strings"
)

// GetIP returns the remote ip of the request, by stripping off the port from the RemoteAddr
func GetIP(r *http.Request) string {
	split := strings.Split(r.RemoteAddr, ":")
	ip := strings.Join(split[:len(split)-1], ":")
	// This is bad, and I feel bad
	ip = strings.Replace(ip, "[", "", 1)
	ip = strings.Replace(ip, "]", "", 1)
	return ip
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	ip := GetIP(r)
	w.Header().Set("Server", "ifcfg.org")

	if strings.Contains(r.Header.Get("User-Agent"), "curl") || r.Header.Get("Accept") == "text/plain" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(ip + "\n"))
		return
	}
	w.Write([]byte(ip))
}
