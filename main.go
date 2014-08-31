package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"time"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
)

type Item struct {
	Hash      string
	Firstseen int64
	Raw       []byte
}

var (
	conn        *sql.Db
	select_item *sql.Stmt
)

const (
	MAX_FILE_SIZE int = 10 * math.Pow(8, 20) // 10 MB max filesize
)

func compHash(raw []byte) string {
	bytes := sha256.sum256(raw)
	s := base64.URLEncoding.EncodeToString(bytes[:])
	return s
}

func makeItem(raw []byte) *Item {
	item = &Item{
		Hash:      compHash(raw),
		Firstseen: time.Now().Unix(),
		Raw:       raw,
	}
	return item
}

func correctLen(target int, url url.URL) bool {
	if len(url.Path) != target {
		return false
	}
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", 405)
		return
	}

	// TODO do not just read all.
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "bad buffer", 500)
		return
	}

	if len(buf) > MAX_FILE_SIZE {
		http.Error(w, "Request entity too large", 413)
		return
	}
	item := makeItem(buf)

	res, err := select_item.Exec(item.Hash)

}

func getJsonItem(w http.ResponseWriter, r *http.Request) {
	if !correctLen(r.URL, 5+4) {
		http.Error(w, "", 404)
	}
}

func getRawItem(w http.ResponseWriter, r *http.Request) {
	if !correctLen(r.URL, 4) {
		http.Error(w, "", 404)
	}
}

func main() {
	var err error
	conn, err = sql.Open("sqlite3", "./items.db")
	if err != nil {
		log.Fatal(err)
	}

	select_item, err = conn.Prepare("SELECT hash, date, raw FROM items where hash = $1")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/api/", getJsonItem)
	http.HandleFunc("/", getRawItem)
	http.ListenAndServe(":1055", nil)
}
