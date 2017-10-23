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

// RootCmd 基命令
var RootCmd *cobra.Command

var git, patch, directory string

func init() {
	var thread int

	transport := http.Transport{
		ResponseHeaderTimeout: time.Second * 20,
		DisableKeepAlives:     true,
		// IdleConnTimeout:       time.Second * 60,
		MaxIdleConns: thread,
	}

	client := http.Client{
		Transport: &transport,
	}

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
			ch := make(chan int, thread)
			if !db.GetBool(base) {
				ch <- 1
				wg.Add(1)
				go spider(base, db, client, ch)
				wg.Wait()
			}

			iter := db.NewIterator(util.BytesPrefix([]byte(tree)), nil)
			for iter.Next() {
				key := string(iter.Key())
				value := string(iter.Value())
				if value == "0" {
					ch <- 1
					wg.Add(1)
					go spider(key, db, client, ch)
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
					// fmt.Println(strings.TrimPrefix(key, tree))
					os.MkdirAll(strings.TrimPrefix(key, tree), os.ModeDir|0755)
				}
			}
			iter.Release()
			if iter.Error() != nil {
				fmt.Println("error")
			}

			ch = make(chan int, thread)
			iter = db.NewIterator(util.BytesPrefix([]byte(raw)), nil)
			for iter.Next() {
				key := string(iter.Key())
				value := string(iter.Value())
				sum++
				if value == "0" {
					fmt.Println(key, "start")
					suc++
					ch <- 1
					wg.Add(1)
					go download(client, key, db, ch)
				}
			}
			iter.Release()
			if iter.Error() != nil {
				fmt.Println("error")
			}
			fmt.Println(suc, "/", sum)
			wg.Wait()
		},
	}
	RootCmd.Flags().StringVarP(&git, "git", "g", "https://gitlab.com/TeeFirefly/FireNow-Nougat/", "git url")
	RootCmd.Flags().StringVarP(&patch, "patch", "p", "master", "patch")
	RootCmd.Flags().StringVarP(&directory, "directory", "d", "", "directories path")
	RootCmd.Flags().IntVarP(&thread, "thread", "t", 100, "thread num")
}

func spider(path string, db *dao.LinkDB, client http.Client, ch chan int) {
	fmt.Println("start ", path)
	batch := new(leveldb.Batch)
	url := `https://gitlab.com` + path

	res, err := client.Get(url)
	for {
		if err == nil {
			break
		}
		fmt.Println(err)
		time.Sleep(time.Second * 3)
		res, err = client.Get(url)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		log.Fatalln(err)
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
					ch <- 1
					wg.Add(1)
					go spider(href, db, client, ch)
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
	<-ch
	wg.Done()
}

func download(client http.Client, path string, db *dao.LinkDB, ch chan int) {
	res, err := client.Get(`https://gitlab.com` + path)

	if err != nil {
		fmt.Println("get", err)
		time.Sleep(time.Second * 3)
		download(client, path, db, ch)
		return
	}
	defer res.Body.Close()

	f, err := os.Create(strings.TrimPrefix(path, raw))
	defer f.Close()
	if err != nil {
		fmt.Println(path, err)
		download(client, path, db, ch)
		return
	}
	len, err := io.Copy(f, res.Body)
	if err != nil {
		fmt.Println(path, "copy", err)
		download(client, path, db, ch)
		return
	}
	fmt.Println(path, len, res.ContentLength)

	if err := db.PutBool(path, true); err != nil {
		fmt.Println(err)
		download(client, path, db, ch)
		return
	}
	<-ch
	wg.Done()
}
