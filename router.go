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
)

var router *mux.Router

func setupRouter() {
	log.Info("Setting up router")
	router = mux.NewRouter()

	slawniakComRouter := router.Host("slawniak.com").PathPrefix("/").Name("slawniak.com").Subrouter()
	slawniakComRouter.PathPrefix("/").HandlerFunc(indexHandler)

	// This seems so wrong
	router.Host("ifcfg.org").Name("ifcfg.org").PathPrefix("/").HandlerFunc(rootHandler)
	router.Host("v4.ifcfg.org").Name("ifcfg.org-v4").PathPrefix("/").HandlerFunc(rootHandler)
	router.Host("v6.ifcfg.org").Name("ifcfg.org-v6").PathPrefix("/").HandlerFunc(rootHandler)

	router.PathPrefix("/").HandlerFunc(indexHandler).Name("catch-all")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if !domainIsRegistered(r.Host) {
		// Make sure we do this syncronousley
		addToDomainList(r.Host)
	}

	staticFolder := "./sites/" + r.Host
	if _, err := os.Stat(staticFolder); err != nil {
		staticFolder = "./client"
	}

	if inf, err := os.Stat(staticFolder + path); err == nil && !inf.IsDir() {
		serveFile(w, r, staticFolder+path)
	} else {
		serveFile(w, r, staticFolder+"/index.html")
	}
}
