package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"lastsweetop.com/dao"
)

var wg sync.WaitGroup

var tree, blob, raw string

func init() {
	project := "/TeeFirefly/FireNow-Nougat/"
	tree = project + "tree"
	blob = project + "blob"
	raw = project + "raw"
}

func main() {
	db := dao.NewLinkDB("db")
	defer db.Close()

	sum := 0
	base := "/TeeFirefly/FireNow-Nougat/tree/firefly-rk3399/frameworks/base"
	db.PutBool(base, false)
	wg.Add(1)
	go spiler(base, db)
	wg.Wait()

	iter := db.NewIterator(util.BytesPrefix([]byte(tree)), nil)
	for iter.Next() {
		key := string(iter.Key())
		value := string(iter.Value())
		if value == "0" {
			wg.Add(1)
			go spiler(key, db)
		} else {
			// fmt.Println("exist", key)
			sum++
			fmt.Println(sum)
			os.MkdirAll(strings.TrimLeft(key, tree), os.ModeDir|0755)
		}
	}
	iter.Release()
	if iter.Error() != nil {
		fmt.Println("error")
	}
	wg.Wait()

	sum = 0
	iter = db.NewIterator(util.BytesPrefix([]byte(raw)), nil)
	for iter.Next() {
		sum++
	}
	fmt.Println(sum)
	iter.Release()
	if iter.Error() != nil {
		fmt.Println("error")
	}
}

func spiler(path string, db *dao.LinkDB) {
	fmt.Println("start ", path)
	batch := new(leveldb.Batch)
	url := `https://gitlab.com` + path
	doc, err := goquery.NewDocument(url)
	for err != nil {
		// fmt.Println(url, "fail load")
		time.Sleep(20000)
		doc, err = goquery.NewDocument(url)
	}
	selection := doc.Find(".tree-item-file-name")
	selection.Each(func(i int, s *goquery.Selection) {
		link := s.Find("a")
		if href, has := link.Attr("href"); has {
			// fmt.Println(link.Text(), `https://gitlab.com`+href)
			if strings.HasPrefix(href, tree) {
				if !strings.HasSuffix(href, "..") && !db.GetBool(href) {
					batch.Put([]byte(href), []byte("0"))
					// db.PutBool(href, false)
					wg.Add(1)
					go spiler(href, db)
				}
			}
			if strings.HasPrefix(href, blob) {
				batch.Put([]byte(strings.Replace(href, blob, raw, -1)), []byte("0"))
			}
		}
	})
	// fmt.Println("success", path)
	// db.PutBool(path, true)
	batch.Put([]byte(path), []byte("1"))
	err = db.Write(batch, nil)
	if err != nil {
		log.Fatalln(err)
	}
	// os.MkdirAll(strings.TrimLeft(path, "/TeeFirefly/FireNow-Nougat/tree/"), os.ModeDir|0755)
	wg.Done()
}
