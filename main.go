package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
)

type Item struct {
	Hash        string
	FirstSeen   int64
	ContentType string
	Raw         []byte
}

var (
	conn          *sql.DB
	selectItem    *sql.Stmt
	storeItem     *sql.Stmt
	MAX_FILE_SIZE int = 10 * int(math.Pow(2, 20)) // 10 MB max filesize
)

func compHash(raw []byte) string {
	bytes := sha256.Sum256(raw)
	s := base64.URLEncoding.EncodeToString(bytes[:])
	return s
}

func makeItem(raw []byte, contType string) *Item {
	item := &Item{
		Hash:        compHash(raw),
		FirstSeen:   time.Now().Unix(),
		ContentType: contType,
		Raw:         raw,
	}
	return item
}

func scanItem(rows *sql.Rows) (*Item, error) {
	var hash, conttype string
	var firstseen int64
	var raw []byte = make([]byte, 0)

	err := rows.Scan(&hash, &firstseen, &conttype, &raw)
	if err != nil {
		return nil, err
	} else {
		item := &Item{
			Hash:        hash,
			FirstSeen:   firstseen,
			ContentType: conttype,
			Raw:         raw,
		}
		return item, nil
	}
}

func correctLen(target int, url *url.URL) bool {
	if len(url.Path) != target {
		return false
	}
	return true
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed.", 405)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Could not read file.", 500)
		return
	}

	cont_type := r.Header.Get("Content-Type")
	if cont_type == "" {
		http.Error(w, "No Content-Type provided", 403)
	}
	item := makeItem(buf, cont_type)

	rows, err := selectItem.Query(item.Hash)
	defer rows.Close()
	if err != nil {
		http.Error(w, "Database error.", 500)
	}
	if rows.Next() {
		http.Error(w, "Entity already stored", 409)
		return
	}

	res, err := storeItem.Exec(item.Hash, item.FirstSeen, item.ContentType, item.Raw)
	if err != nil {
		http.Error(w, "Database error.", 500)
		return
	}
	i, err := res.RowsAffected()
	if i != 1 || err != nil {
		http.Error(w, "Database error.", 500)
		return
	}

	itemUrl := "http://localhost:1055/" + item.Hash
	w.WriteHeader(201)
	w.Write([]byte(itemUrl))
}

func getItem(w http.ResponseWriter, itemHash string) *Item {
	rows, err := selectItem.Query(itemHash)
	defer rows.Close()
	if err != nil {
		http.Error(w, "Database error.", 500)
		return nil
	}
	if !rows.Next() {
		http.Error(w, "Item not found in database.", 404)
		return nil
	}

	item, err := scanItem(rows)
	if err != nil {
		http.Error(w, "Database error.", 500)
		return nil
	}
	return item
}

func getJsonItem(w http.ResponseWriter, r *http.Request) {
	if !correctLen(len("/api/")+44, r.URL) {
		http.Error(w, "Improper url format", 500)
	}
	itemHash := r.URL.Path[5:]
	item := getItem(w, itemHash)
	if item == nil {
		return
	}

	bytes, err := json.Marshal(item)
	if err != nil {
		http.Error(w, "JSON marshal failed.", 500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(item.Raw)))
	w.Write(bytes)
}

func getRawItem(w http.ResponseWriter, r *http.Request) {
	if !correctLen(45, r.URL) {
		http.Error(w, "Improper url format", 500)
		return
	}

	itemHash := r.URL.Path[1:]
	item := getItem(w, itemHash)
	if item == nil {
		return
	}

	w.Header().Set("Content-Type", item.ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(item.Raw)))
	w.Write(item.Raw)
}

func main() {
	var err error
	conn, err = sql.Open("sqlite3", "./items.db")
	if err != nil {
		log.Fatal(err)
	}

	selectItem, err = conn.Prepare("SELECT hash, firstseen, conttype, raw FROM items where hash = $1")
	if err != nil {
		log.Fatal(err)
	}
	storeItem, err = conn.Prepare("INSERT INTO items (hash, firstseen, conttype, raw) VALUES($1, $2, $3, $4)")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/api/", getJsonItem)
	http.HandleFunc("/", getRawItem)
	log.Fatal(http.ListenAndServe(":1055", nil))
}
