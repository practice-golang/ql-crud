package main

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/cznic/ql"
)

// Book : 책 정보
type Book struct {
	ID     uint
	Title  string
	Author string
}

func createTable(table string, db *ql.DB) error {
	_, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		create table if not exists `+table+` (id bigint, title string, author string);
		commit;`)
	if err != nil {
		return err
	}

	return nil
}

func insertData(book *Book, table string, db *ql.DB) error {
	s, _, err := db.Run(ql.NewRWCtx(), `begin transaction; select id() from `+table+` order by id() desc limit 1; commit;`)
	if err != nil {
		return err
	}
	lastID := uint64(1)
	lastData, _ := s[0].FirstRow()
	if len(lastData) > 0 {
		lastID = uint64(lastData[0].(int64) + 1)
	}

	_, _, err = db.Run(ql.NewRWCtx(),
		`begin transaction;
		insert into `+table+` values (`+strconv.FormatUint(lastID, 10)+`, "`+book.Title+`", "`+book.Author+`");
		commit;`)
	if err != nil {
		return err
	}

	return nil
}

func updateData(book *Book, table string, db *ql.DB) error {
	idStr := strconv.FormatUint(uint64(book.ID), 16)

	_, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		update `+table+` set title="`+book.Title+`", author="`+book.Author+`" where id==`+idStr+`;
		commit;`)
	if err != nil {
		return err
	}

	return nil
}

func deleteData(id uint64, db *ql.DB) error {
	idStr := strconv.FormatUint(id, 16)

	_, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		delete from books where id==`+idStr+`;
		commit;`)
	if err != nil {
		return err
	}

	return nil
}

func selectData(table string, db *ql.DB) ([][]interface{}, error) {
	s, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		select id, title, author from `+table+` order by id() desc;
		commit;`)
	if err != nil {
		return nil, err
	}

	data, _ := s[0].Rows(99, 0)
	// fmt.Println(data)
	// fmt.Println(reflect.TypeOf(data))

	return data, nil
}

func dropTable(table string, db *ql.DB) error {
	_, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		drop table `+table+`;
		commit;`)
	if err != nil {
		return err
	}

	return nil
}

const dbname = "./ql.db"

func main() {
	db, err := ql.OpenFile(dbname, &ql.Options{CanCreate: true, RemoveEmptyWAL: true})
	defer db.Close()
	if err != nil {
		fmt.Println("DB Exist", err)
	}

	table := "books"
	book := Book{Title: "First Book", Author: "Steve Rogers"}

	if err = createTable(table, db); err != nil {
		fmt.Println("Table Exist", err)
	}

	if err = insertData(&book, table, db); err != nil {
		panic(err)
	}

	data, err := selectData(table, db)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)
	for _, v := range data {
		book.ID = uint(v[0].(*big.Int).Uint64())
		book.Title = v[1].(string)
		book.Author = v[2].(string)
	}

	book.Title = "First Avenger"

	if err = updateData(&book, table, db); err != nil {
		panic(err)
	}

	data, err = selectData(table, db)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	if err = dropTable(table, db); err != nil {
		panic(err)
	}
}
