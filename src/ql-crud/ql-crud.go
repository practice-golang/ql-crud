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

// CreateTable : 테이블 생성
func CreateTable(table string, db *ql.DB) error {
	_, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		create table if not exists `+table+` (id bigint, title string, author string);
		commit;`)
	if err != nil {
		return err
	}

	return nil
}

// DropTable : 테이블 드랍
func DropTable(table string, db *ql.DB) error {
	_, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		drop table `+table+`;
		commit;`)
	if err != nil {
		return err
	}

	return nil
}

// InsertData : Crud
func InsertData(book *Book, table string, db *ql.DB) error {
	s, _, err := db.Run(ql.NewRWCtx(), `begin transaction; select id() from `+table+` order by id() desc limit 1; commit;`)
	if err != nil {
		panic(err)
		// return err
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
		panic(err)
		// return err
	}

	return nil
}

// SelectData : Crud
func SelectData(table string, db *ql.DB) ([][]interface{}, error) {
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

// UpdateData : crUd
func UpdateData(book *Book, table string, db *ql.DB) error {
	idStr := strconv.FormatUint(uint64(book.ID), 16)

	_, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		update `+table+` set title="`+book.Title+`", author="`+book.Author+`" where id==`+idStr+`;
		commit;`)
	if err != nil {
		panic(err)
		// return err
	}

	return nil
}

// DeleteData : cruD
func DeleteData(id uint64, db *ql.DB) error {
	idStr := strconv.FormatUint(id, 16)

	_, _, err := db.Run(ql.NewRWCtx(),
		`begin transaction;
		delete from books where id==`+idStr+`;
		commit;`)
	if err != nil {
		panic(err)
		// return err
	}

	return nil
}

const dbname = "./ql.db"

func main() {
	table := "books"
	var books []Book
	book := Book{Title: "First Book", Author: "Steve Rogers"}

	// DB 준비. 없으면 생성. ql.Open은 lock 파일 삭제 방도를 못 찾아서 일단 안 씀.
	db, err := ql.OpenFile(dbname, &ql.Options{CanCreate: true, RemoveEmptyWAL: true})
	defer db.Close()
	if err != nil {
		fmt.Println("DB Exist", err)
	}

	if err = CreateTable(table, db); err != nil {
		fmt.Println("Table Exist", err)
	}

	InsertData(&book, table, db) // Crud

	data, _ := SelectData(table, db) // cRud
	fmt.Println(data)

	// 배열이나 슬라이스에서 구조체로 예쁘게 가져올 방법 없나?
	for _, v := range data {
		book.ID, book.Title, book.Author = uint(v[0].(*big.Int).Uint64()), v[1].(string), v[2].(string)
		books = append(books, book)
	}

	books[0].Title = "First Avenger"
	UpdateData(&books[0], table, db) // crUd

	data, _ = SelectData(table, db)
	fmt.Println(data)

	// 테이블을 드랍해도 ql에서 쓰이는 id()는 계속 갱신되는 것 같다.
	if err = DropTable(table, db); err != nil {
		panic(err)
	}
}
