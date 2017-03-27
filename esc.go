package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
)

type _escLocalFS struct{}

var _escLocal _escLocalFS

type _escStaticFS struct{}

var _escStatic _escStaticFS

type _escDirectory struct {
	fs   http.FileSystem
	name string
}

type _escFile struct {
	compressed string
	size       int64
	modtime    int64
	local      string
	isDir      bool

	once sync.Once
	data []byte
	name string
}

func (_escLocalFS) Open(name string) (http.File, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	return os.Open(f.local)
}

func (_escStaticFS) prepare(name string) (*_escFile, error) {
	f, present := _escData[path.Clean(name)]
	if !present {
		return nil, os.ErrNotExist
	}
	var err error
	f.once.Do(func() {
		f.name = path.Base(name)
		if f.size == 0 {
			return
		}
		var gr *gzip.Reader
		b64 := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(f.compressed))
		gr, err = gzip.NewReader(b64)
		if err != nil {
			return
		}
		f.data, err = ioutil.ReadAll(gr)
	})
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (fs _escStaticFS) Open(name string) (http.File, error) {
	f, err := fs.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.File()
}

func (dir _escDirectory) Open(name string) (http.File, error) {
	return dir.fs.Open(dir.name + name)
}

func (f *_escFile) File() (http.File, error) {
	type httpFile struct {
		*bytes.Reader
		*_escFile
	}
	return &httpFile{
		Reader:   bytes.NewReader(f.data),
		_escFile: f,
	}, nil
}

func (f *_escFile) Close() error {
	return nil
}

func (f *_escFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (f *_escFile) Stat() (os.FileInfo, error) {
	return f, nil
}

func (f *_escFile) Name() string {
	return f.name
}

func (f *_escFile) Size() int64 {
	return f.size
}

func (f *_escFile) Mode() os.FileMode {
	return 0
}

func (f *_escFile) ModTime() time.Time {
	return time.Unix(f.modtime, 0)
}

func (f *_escFile) IsDir() bool {
	return f.isDir
}

func (f *_escFile) Sys() interface{} {
	return f
}

// FS returns a http.Filesystem for the embedded assets. If useLocal is true,
// the filesystem's contents are instead used.
func FS(useLocal bool) http.FileSystem {
	if useLocal {
		return _escLocal
	}
	return _escStatic
}

// Dir returns a http.Filesystem for the embedded assets on a given prefix dir.
// If useLocal is true, the filesystem's contents are instead used.
func Dir(useLocal bool, name string) http.FileSystem {
	if useLocal {
		return _escDirectory{fs: _escLocal, name: name}
	}
	return _escDirectory{fs: _escStatic, name: name}
}

// FSByte returns the named file from the embedded assets. If useLocal is
// true, the filesystem's contents are instead used.
func FSByte(useLocal bool, name string) ([]byte, error) {
	if useLocal {
		f, err := _escLocal.Open(name)
		if err != nil {
			return nil, err
		}
		b, err := ioutil.ReadAll(f)
		f.Close()
		return b, err
	}
	f, err := _escStatic.prepare(name)
	if err != nil {
		return nil, err
	}
	return f.data, nil
}

// FSMustByte is the same as FSByte, but panics if name is not present.
func FSMustByte(useLocal bool, name string) []byte {
	b, err := FSByte(useLocal, name)
	if err != nil {
		panic(err)
	}
	return b
}

// FSString is the string version of FSByte.
func FSString(useLocal bool, name string) (string, error) {
	b, err := FSByte(useLocal, name)
	return string(b), err
}

// FSMustString is the string version of FSMustByte.
func FSMustString(useLocal bool, name string) string {
	return string(FSMustByte(useLocal, name))
}

var _escData = map[string]*_escFile{

	"/index.html": {
		local:   "client/index.html",
		size:    420,
		modtime: 1490582055,
		compressed: `
H4sIAAAJbogA/2SRP0/DMBDF90j5DuaYk0ggEIOdodCuMLAwGvtSm/pPZR+N8u1REkMrOtl+z/d+Tza/
eXl9fv942zJD3vV1xeeVORn2AjBAX1ezhlL3dcUY90iSKSNTRhLwTUPzBBdOkB4FnCyOx5gImIqBMJCA
0WoyQuPJKmyWQxlzNhxYQifAy2AHzATMJBwEdL9C+5VjANZdccigx0ZFF9MF6na3227uN8CsLlf+Zs8w
q+bIAsokyapu1pq7h8f2GPalHVly2BsMaWqzk2Ow8tCq6Hm3OudUjQOmNTvT5DAbRPpPWJxW5bw24l15
2Lrin1FPi1Q2i7t+yU8AAAD//ycHyMikAQAA
`,
	},

	"/manifest.json": {
		local:   "client/manifest.json",
		size:    325,
		modtime: 1490580139,
		compressed: `
H4sIAAAJbogA/4yO0QqCMBSG7/cUY92W5shd+CoRstbQodtkm6SF7x6/EnQRFAd2cb7/33eehFLmpNWs
oqzVLsxZ7OXdGdllylu2B49JhlSPoUco33ZXqbom+NHdauV7H4B2vMBsgdRqqz+ZwGzMKO8iq+gZegiC
Wr+OSSajcuADL0U2uGYtIGIeGhXGSzHxUrz3aR7W242Vjc5RIJQul781xZGfvnoAJjw/TGQhrwAAAP//
BwQAc0UBAAA=
`,
	},

	"/static/icon-1024.png": {
		local:   "client/static/icon-1024.png",
		size:    16919,
		modtime: 1490581396,
		compressed: `
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
/8/91uYXQgAA
`,
	},

	"/static/icon-256.png": {
		local:   "client/static/icon-256.png",
		size:    4364,
		modtime: 1490581376,
		compressed: `
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
BQAAjw1b3SvdQo7+TwAAAP//ab3u3wwRAAA=
`,
	},

	"/static/style.css": {
		local:   "client/static/style.css",
		size:    53,
		modtime: 1490580161,
		compressed: `
H4sIAAAJbogA/8ooyc3RUUjKT6lUqOblUlDIzczTzUjNTM8osVIwNDAoy7AGiyYWpWfmWSkYWPNy1fJy
AQIAAP//idtB8DUAAAA=
`,
	},

	"/": {
		isDir: true,
		local: "client",
	},

	"/static": {
		isDir: true,
		local: "client/static",
	},
}
