package download

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/mholt/archiver"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
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

	f, e := os.CreateTemp(dstPath, "ipdb-")
	if e != nil {
		return e
	}
	_, e = io.Copy(f, res.Body)
	if e != nil {
		return e
	}
	fn := f.Name()
	f.Close()

	all, e := os.ReadFile(fn)
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

	var dst string

	// : attachment; filename="idc_list.ipdb"
	g := regexp.MustCompile(`filename="([^"]+)"`).FindAllStringSubmatch(res.Header.Get("Content-Disposition"), -1)
	if len(g) < 1 {
		fmt.Println(g)
		return fmt.Errorf("download attachment failed")
	}

	name = g[0][1]

	if strings.HasSuffix(name, ".zip") {
		Z := archiver.NewZip()
		e = Z.Walk(fn, func(f archiver.File) error {

			if f.IsDir() {
				return nil
			}
			defer f.Close()
			dst = filepath.Join(dstPath, f.Name())
			w, e := os.Create(dst)
			if e != nil {
				log.Println(e)
				return e
			}
			defer w.Close()
			_, e = io.Copy(w, f)
			if e != nil {
				log.Println(e)
				return e
			}
			return nil
		})

		os.Remove(fn) // 删除临时zip包
	} else {
		return os.Rename(fn, filepath.Join(dstPath, name))
	}

	return e
}
