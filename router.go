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
	"github.com/go-playground/log"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strings"
)

var router *mux.Router

func setupRouter() {
	log.Info("Setting up router")
	router = mux.NewRouter()

	slawniakComRouter := router.Host("slawniak.com").PathPrefix("/").Name("slawniak.com").Subrouter()
	slawniakComRouter.PathPrefix("/").HandlerFunc(indexHandler)

	// This seems so wrong
	router.Host("ifcfg.org").Name("ifcfg.org").PathPrefix("/").HandlerFunc(ifcfgRootHandler)
	router.Host("v4.ifcfg.org").Name("ifcfg.org-v4").PathPrefix("/").HandlerFunc(ifcfgRootHandler)
	router.Host("v6.ifcfg.org").Name("ifcfg.org-v6").PathPrefix("/").HandlerFunc(ifcfgRootHandler)

	router.Host("stopallthe.download").Path("/ing/provision").Handler(http.RedirectHandler("https://gist.githubusercontent.com/HenrySlawniak/c31cedaec491c68631a6f62b5d94a740/raw", http.StatusFound))
	router.Host("stopallthe.download").Path("/ing/install-go").Handler(http.RedirectHandler("https://gist.githubusercontent.com/HenrySlawniak/1b17dc248f57016ee820a7502d7285ce/raw", http.StatusFound))
	router.Host("stopallthe.download").Name("stopall").PathPrefix("/ing/").HandlerFunc(stopAllIngHandler)
	router.Host("stopallthe.download").Name("stopall").PathPrefix("/").HandlerFunc(stopAllRootHandler)

	router.PathPrefix("/").HandlerFunc(indexHandler).Name("catch-all")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	host := strings.Split(r.Host, ":")[0]

	if !domainIsRegistered(host) {
		// Make sure we do this syncronousley
		addToDomainList(host, true)
	}

	staticFolder := "./sites/" + host
	if _, err := os.Stat(staticFolder); err != nil {
		staticFolder = "./client"
	}

	var n int64
	var code int

	if inf, err := os.Stat(staticFolder + path); err == nil && !inf.IsDir() {
		n, code = serveFile(w, r, staticFolder+path)
	} else if inf, err := os.Stat(staticFolder + path + "/index.html"); err == nil && !inf.IsDir() {
		n, code = serveFile(w, r, staticFolder+path+"/index.html")
	} else {
		n, code = serveFile(w, r, staticFolder+"/index.html")
	}

	go logRequest(w, r, n, code)
}
