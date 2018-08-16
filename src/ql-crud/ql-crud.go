package main

import (
	"fmt"

	"github.com/cznic/ql"
)

func createTable(db *ql.DB) error {
	_, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		create table if not exists books (id bigint, title string, author string);
		commit;`)
	if err != nil {
		return err
	}

	return nil
}

func insertData(db *ql.DB) error {
	_, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		insert into books (id, title, author) values (1, "Book Title", "Author Man");
		insert into books (id, title, author) values (2, "Titleaaa", "Girl");
		insert into books (title, author) values ("Booooooooooook", "Childddd");
		commit;`)
	if err != nil {
		return err
	}

	return nil
}

func selectData(db *ql.DB) error {
	a, b, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		select * from books order by id asc;
		commit;`)
	if err != nil {
		return err
	}

	fmt.Println(a)
	data, _ := a[0].Rows(99, 0)
	fmt.Println(data)
	fmt.Println(b)

	return nil
}

func dropTable(db *ql.DB) error {
	_, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		drop table books;
		commit;`)
	if err != nil {
		return err
	}
	return nil
}

const dbname = "ql.db"

func main() {
	db, err := ql.OpenFile(dbname, &ql.Options{CanCreate: true, RemoveEmptyWAL: true})
	defer db.Close()
	if err != nil {
		panic(err)
	}

	if err = createTable(db); err != nil {
		panic(err)
	}
	if err = insertData(db); err != nil {
		panic(err)
	}
	if err = selectData(db); err != nil {
		panic(err)
	}
	if err = dropTable(db); err != nil {
		panic(err)
	}
}
