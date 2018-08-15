package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/cznic/ql"
)

func setUp(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`create table books (id bigint , title string, body string, created_at string, updated_at string);`)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func tearDown(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`drop table note;`)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

const dbname = "./ql.db"

func walName(dbname string) (r string) {
	base := filepath.Base(filepath.Clean(dbname))
	h := sha1.New()
	io.WriteString(h, base)

	return filepath.Join(filepath.Dir(dbname), fmt.Sprintf(".%x", h.Sum(nil)))
}

func main() {
	wName := walName(dbname)

	ql.RegisterDriver()

	db, err := sql.Open("ql", dbname)
	defer os.Remove(wName)
	defer db.Close()
	if err != nil {
		log.Fatalf("failed to open db: %s", err)
	}

	if err = setUp(db); err != nil {
		log.Fatalf("failed to create table: %s", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %s", err)
	}
	if err = tearDown(db); err != nil {
		log.Fatalf("failed to drop table: %s", err)
	}
}
