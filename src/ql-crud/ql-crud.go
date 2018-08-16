package main

import (
	"fmt"
	"strconv"

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
		commit;`)
	if err != nil {
		return err
	}

	a, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		select id() from books order by id() desc limit 1;
		commit;`)
	if err != nil {
		return err
	}
	lastData, _ := a[0].FirstRow()
	lastID := uint64(lastData[0].(int64) + 1)

	_, _, err = db.Run(ql.NewRWCtx(),
		`begin transaction;
		insert into books values (`+strconv.FormatUint(lastID, 10)+`, "Booooooooooook", "Childddd");
		commit;`)
	if err != nil {
		fmt.Println("WTF::", err)
		return err
	}

	return nil
}

func updateData(db *ql.DB) error {
	_, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		update books set title="Titan^_^_^", author="Woman" where id==2;
		commit;`)
	if err != nil {
		return err
	}

	return nil
}

func selectData(db *ql.DB) error {
	a, b, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		select id(), id, title, author from books order by id() desc;
		commit;`)
	if err != nil {
		return err
	}

	fmt.Println(a)
	data, _ := a[0].Rows(999, 0)
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
	if err = updateData(db); err != nil {
		panic(err)
	}
	if err = selectData(db); err != nil {
		panic(err)
	}
	if err = dropTable(db); err != nil {
		panic(err)
	}
}
