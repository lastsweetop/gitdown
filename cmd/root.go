package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lastsweetop/gitdown/dao"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var wg sync.WaitGroup

var tree, blob, raw string
var RootCmd *cobra.Command

func init() {
	var git, patch, directory string

	RootCmd = &cobra.Command{
		Use: "gitdown",
		Run: func(cmd *cobra.Command, args []string) {

			fmt.Println(git)
			project := strings.TrimSuffix(strings.TrimPrefix(git, `https://gitlab.com`), ".git")
			fmt.Println(project)
			tree = project + "/tree/"
			blob = project + "/blob/"
			raw = project + "/raw/"

			db := dao.NewLinkDB("db")
			defer db.Close()

			sum := 0
			suc := 0
			base := project + "/tree/" + patch + "/" + directory

			db.PutBool(base, false)
			wg.Add(1)
			go spider(base, db)
			wg.Wait()

			iter := db.NewIterator(util.BytesPrefix([]byte(tree)), nil)
			for iter.Next() {
				key := string(iter.Key())
				value := string(iter.Value())
				if value == "0" {
					wg.Add(1)
					go spider(key, db)
				}
			}
			iter.Release()
			if iter.Error() != nil {
				fmt.Println("error")
			}
			wg.Wait()

			iter = db.NewIterator(util.BytesPrefix([]byte(tree)), nil)
			for iter.Next() {
				key := string(iter.Key())
				value := string(iter.Value())
				if value == "1" {
					fmt.Println(strings.TrimPrefix(key, tree))
					os.MkdirAll(strings.TrimPrefix(key, tree), os.ModeDir|0755)
				}
			}
			iter.Release()
			if iter.Error() != nil {
				fmt.Println("error")
			}

			iter = db.NewIterator(util.BytesPrefix([]byte(raw)), nil)
			for iter.Next() {
				key := string(iter.Key())
				value := string(iter.Value())
				sum++
				if value == "0" {
					suc++
					wg.Add(1)
					go download(key, db)
				}
			}
			fmt.Println(suc, "/", sum)
			iter.Release()
			if iter.Error() != nil {
				fmt.Println("error")
			}
			wg.Wait()
		},
	}
	RootCmd.Flags().StringVarP(&git, "git", "g", "https://gitlab.com/TeeFirefly/FireNow-Nougat/", "git url")
	RootCmd.Flags().StringVarP(&patch, "patch", "p", "master", "patch")
	RootCmd.Flags().StringVarP(&directory, "directory", "d", "", "directories path")
}

func spider(path string, db *dao.LinkDB) {
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
			// fmt.Println(link.Text(), `https://gitlab.com`+href, tree)
			if strings.HasPrefix(href, tree) {
				if !strings.HasSuffix(href, "..") && !db.GetBool(href) {
					batch.Put([]byte(href), []byte("0"))
					// db.PutBool(href, false)
					wg.Add(1)
					go spider(href, db)
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

func download(path string, db *dao.LinkDB) {
	res, err := http.Get(`https://gitlab.com` + path)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(strings.TrimPrefix(path, raw))
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	len, err := io.Copy(f, res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(path, len)
	if err := db.PutBool(path, true); err != nil {
		log.Fatal(err)
	}
	wg.Done()
}
