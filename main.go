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
	"bufio"
	"crypto/tls"
	"flag"
	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	devMode            = flag.Bool("dev", false, "Puts the server in developer mode, will bind to :34265 and will not autocert")
	accessLogInConsole = flag.Bool("console-access", false, "Whether or not to print access log lines to the console")
	cookieSecret       string
	buildTime          string
	commit             string
	domainList         = []string{}
	m                  autocert.Manager
)

func init() {
	flag.Parse()
	cLog := console.New()
	cLog.SetTimestampFormat(time.RFC3339)
	log.RegisterHandler(cLog, log.AllLevels...)
}

func main() {
	log.Info("Starting henry.slawniak.com")
	if buildTime != "" {
		log.Info("Built: " + buildTime)
	}
	if commit != "" {
		log.Info("Revision: " + commit)
	}
	log.Info("Go: " + runtime.Version())
	setupRouter()

	loadDomainList()

	if *devMode {
		srv := &http.Server{
			Addr:    ":34265",
			Handler: router,
		}

		log.Info("Listening on :34265")
		srv.ListenAndServe()
	}

	m = autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domainList...),
		Cache:      autocert.DirCache("certs"),
	}

	httpSrv := &http.Server{
		Addr:    ":http",
		Handler: http.HandlerFunc(httpRedirectHandler),
	}

	go httpSrv.ListenAndServe()

	rootSrv := &http.Server{
		Addr:      ":https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   router,
	}

	log.Info("Listening on :https")

	http2.ConfigureServer(rootSrv, &http2.Server{})
	rootSrv.ListenAndServeTLS("", "")
}

func httpRedirectHandler(w http.ResponseWriter, r *http.Request) {
	if !domainIsRegistered(r.Host) {
		addToDomainList(r.Host, true)
	}

	http.Redirect(w, r, "https://"+r.Host+r.URL.String(), http.StatusMovedPermanently)
}

func domainIsRegistered(domain string) bool {
	for _, d := range domainList {
		if d == domain {
			return true
		}
	}
	return false
}

func addToDomainList(domain string, isNew bool) {
	if domain == "" {
		log.Warn("Cannot use an empty string as a domain")
		return
	}
	for _, d := range domainList {
		if d == domain {
			log.Noticef("%s already in domain list, returning\n", domain)
			return
		}
	}
	domainList = append(domainList, domain)

	err := ioutil.WriteFile("domains.txt", []byte(strings.Join(domainList, "\n")), os.ModeAppend)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	m.HostPolicy = autocert.HostWhitelist(domainList...)
	if isNew {
		log.Noticef("Added %s to registered domains", domain)
	}
}

func loadDomainList() {
	var f *os.File
	var err error
	if _, err = os.Stat("domains.txt"); os.IsNotExist(err) {
		f, err = os.Create("domains.txt")
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
	} else {
		f, err = os.Open("domains.txt")
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	for scanner.Scan() {
		addToDomainList(scanner.Text(), false)
	}

	log.Noticef("There are now %d domains registered\n", len(domainList))
}
