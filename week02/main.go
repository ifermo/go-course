package main

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

var (
	NotFound = errors.New("data not found")
)

func main() {
	dao := &Dao{
		db: initDB(),
	}
	art, err := dao.QueryArticleById(100)
	if err != nil {
		if errors.Is(err, NotFound) {
			fmt.Println("article not found")
		} else {
			fmt.Println(err)
		}
		return
	}
	fmt.Println(art)
}

type Article struct {
	Id      int64
	Content string
	Author  string
}

func (art *Article) String() string {
	return fmt.Sprintf("Article { id:%d, content:%s, author:%s }", art.Id, art.Content, art.Author)
}

type Dao struct {
	db *sql.DB
}

func (dao *Dao) QueryArticleById(id int64) (*Article, error) {
	art := &Article{}
	row := dao.db.QueryRow("select id,content,author from article where id=?", id)
	err := row.Scan(&art.Id, &art.Content, &art.Author)
	if err == sql.ErrNoRows {
		return nil, errors.Wrap(NotFound, err.Error())
	} else if err != nil {
		return nil, errors.Wrap(errors.New("unknown error"), err.Error())
	}
	return art, nil
}

func initDB() *sql.DB {
	db, err := sql.Open("mysql", "user:password@/dbname")
	if err != nil {
		panic(err)
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}
