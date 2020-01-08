package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"sync"
)

var db *sql.DB
var once sync.Once

func init(){
	GetDataBase()
}

func GetDataBase() *sql.DB {

	once.Do(func() {

		dbinf := fmt.Sprintf("user=%s password=%s dbname=%s host=127.0.0.1 port=5432 sslmode=disable", "kexibq", "kexibq", "forumdb")
		var err error
		db, err = sql.Open("postgres", dbinf)
		if err != nil {
			log.Println("Can't connect to database", err)
		}
		err = db.Ping()
		if err != nil {
			log.Println("error in ping", err)
		}
	})
	return db
}

func CloseDB() {
	db.Close()
}

func Init(filename string) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Cant read file", err)
		return
	}
	str := string(bs)
	_, err = db.Exec(str)
	if err != nil {
		log.Println("Error while db Init Executing script", err)
		return
	}

}