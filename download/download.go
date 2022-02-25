package download

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/mholt/archiver"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Custom(token string, name string, dstPath string) error {

	API := url.URL{
		Scheme:     "https",
		Host:       "user.ipip.net",
		Path:       "/download.php",
	}
	q := API.Query()
	q.Add("a", "custom")
	q.Add("token", token)
	q.Add("lang", "CN")
	q.Add("name", name)

	API.RawQuery = q.Encode()

	req, e := http.NewRequest(http.MethodGet, API.String(), nil)
	if e != nil {
		return e
	}
	req.Header.Set("User-Agent", "IPIP.NET")
	res, e := http.DefaultClient.Do(req)
	if e != nil {
		return e
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		all, e := io.ReadAll(res.Body)
		if e == nil {
			return fmt.Errorf("%s", string(all))
		}
	}

	fn := filepath.Join(os.TempDir(), name)
	f, e := os.Create(fn)
	if e != nil {
		return e
	}
	_, e = io.Copy(f, res.Body)
	if e != nil {
		return e
	}
	f.Close()

	all, e := ioutil.ReadFile(fn)
	if e != nil {
		return e
	}

	if len(all) < 128 {
		return fmt.Errorf("file exception")
	}

	h := sha1.New()
	h.Write(all)
	if !strings.HasSuffix(res.Header.Get("ETag"), hex.EncodeToString(h.Sum(nil))) {
		return fmt.Errorf("ETag diff")
	}

	Z := archiver.NewZip()
	var dst string
	e = Z.Walk(fn, func(f archiver.File) error {

		if f.IsDir() {
			return nil
		}
		defer f.Close()
		dst = filepath.Join(dstPath, f.Name())
		fmt.Println(dst)
		w, e := os.Create(dst)
		if e != nil {
			return e
		}
		defer w.Close()
		_, e = io.Copy(w, f)
		if e != nil {
			return e
		}
		return nil
	})

	return e
}
