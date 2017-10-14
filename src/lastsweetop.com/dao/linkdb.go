package dao

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type LinkDB struct {
	*leveldb.DB
}

func NewLinkDB(dbname string) *LinkDB {
	db, err := leveldb.OpenFile(dbname, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &LinkDB{db}
}

func (db *LinkDB) PutBool(url string, done bool) {
	num := "0"
	if done {
		num = "1"
	}
	db.Put([]byte(url), []byte(num), nil)
}

func (db *LinkDB) GetBool(url string) bool {
	ret, err := db.Has([]byte(url), nil)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	if !ret {
		return false
	}

	temp, err := db.Get([]byte(url), nil)
	if err != nil {
		log.Fatalln(err)
		return false
	}
	b := string(temp)
	if b == "0" {
		return false
	}
	return true
}
