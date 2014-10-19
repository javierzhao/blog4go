package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Dao struct {
	db *sql.DB
}

var Cdao *Dao = nil

func Init() (*Dao, error) {
	if Cdao == nil {
		fmt.Println("Init sql.DB .....")
		dao := new(Dao)
		db, err := sql.Open("mysql", "root:1@/blog4go?charset=utf8")
		CheckErr(err)
		dao.db = db
		Cdao = dao
	}
	return Cdao, nil
}

func (this *Dao) Close() {
	fmt.Println("closeing")
	this.db.Close()
	fmt.Println("closed")
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
