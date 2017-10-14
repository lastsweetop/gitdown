package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"lastsweetop.com/dao"
)

var wg sync.WaitGroup

var db *dao.LinkDB
var pd *dao.LinkDB

func main() {
	db := dao.NewLinkDB("db")
	pd := dao.NewLinkDB("pb")

	defer db.Close()
	defer pd.Close()
	sum := 0
	base := "/TeeFirefly/FireNow-Nougat/tree/firefly-rk3399/frameworks/base"
	db.PutBool(base, false)
	wg.Add(1)
	go spiler(base, db, pd, sum)
	wg.Wait()

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := string(iter.Key())
		value := string(iter.Value())
		if value == "0" {
			wg.Add(1)
			go spiler(key, db, pd, sum)
		} else {
			// fmt.Println("exist", key)
			sum++
			fmt.Println(sum)
			os.MkdirAll(strings.TrimLeft(key, "/TeeFirefly/FireNow-Nougat/tree/"), os.ModeDir|0755)
		}
	}
	iter.Release()
	if iter.Error() != nil {
		fmt.Println("error")
	}
	wg.Wait()

	sum = 0
	iter = pd.NewIterator(nil, nil)
	for iter.Next() {
		sum++
	}
	fmt.Println(sum)
	iter.Release()
	if iter.Error() != nil {
		fmt.Println("error")
	}
}

func spiler(path string, db *dao.LinkDB, pd *dao.LinkDB, sum int) {
	// fmt.Println("start ", path)
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
			if strings.HasPrefix(href, "/TeeFirefly/FireNow-Nougat/tree") {
				if !strings.HasSuffix(href, "..") && !db.GetBool(href) {
					db.PutBool(href, false)
					wg.Add(1)
					go spiler(href, db, pd, sum)
				}
			}
			if strings.HasPrefix(href, "/TeeFirefly/FireNow-Nougat/blob") {
				pd.PutBool(strings.Replace(href, "blob", "raw", -1), false)
			}
		}
	})
	// fmt.Println("success", path)
	db.PutBool(path, true)
	// os.MkdirAll(strings.TrimLeft(path, "/TeeFirefly/FireNow-Nougat/tree/"), os.ModeDir|0755)
	wg.Done()
}
