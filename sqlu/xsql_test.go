package sqlu

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

type A struct {
	School  string
	Address string
	Door    int
}

type Test struct {
	Name  string
	Age   uint
	Socre float32
	A1    A //默认json
	A
}

var p = Test{
	Name:  "John",
	Age:   26,
	Socre: 89.5,
	A1: A{
		School:  "中国医科大",
		Address: "河南南阳",
	},
	A: A{
		School: "世界第一",
	},
}

func TestStmt(t *testing.T) {
	fmt.Println(SqlValueFmt(&p))
	fmt.Println(SqlValues(p, nil))
	var p1 []Test
	for i := 0; i < 10; i++ {
		p1 = append(p1, p)
	}
	fmt.Println(SqlValues(p1, nil))
}

func TestDB(t *testing.T) {
	db, err := sql.Open("mysql", "root:123456@tcp(172.16.60.219:35200)/test")
	if err != nil {
		log.Fatalln(err)
	}
	myDB := &DB{db: db}
	_, err = myDB.Insert(&p, "test")
	if err != nil {
		log.Fatal(err)
	}
}

func TestSqlc(t *testing.T) {
	db, err := sql.Open("mysql", "root:123456@tcp(172.16.60.219:35200)/test")
	if err != nil {
		log.Fatalln(err)
	}
	myDB := &DB{db: db}
	var total int
	sqlc := NewSQL("SELECT count(1) FROM test").Equal("Name", "").Equal("Age", 16)
	if err := myDB.Count(&total, sqlc.String()); err != nil {
		log.Fatal(err)
		return
	}
	log.Println("total", total)
}
