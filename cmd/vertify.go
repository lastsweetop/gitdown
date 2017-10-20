package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lastsweetop/gitdown/dao"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func init() {
	vertifyCmd.Flags().StringVarP(&git, "git", "g", "https://gitlab.com/TeeFirefly/FireNow-Nougat", "git url")
	vertifyCmd.Flags().StringVarP(&patch, "patch", "p", "firefly-rk3399", "patch")
	RootCmd.AddCommand(vertifyCmd)
}

var vertifyCmd = &cobra.Command{
	Use:   "vertify",
	Short: "vertify the repo",
	Run: func(cmd *cobra.Command, args []string) {
		transport := http.Transport{
			ResponseHeaderTimeout: time.Second * 20,
			DisableKeepAlives:     false,
			MaxIdleConns:          100,
		}

		client := http.Client{
			Transport: &transport,
		}
		project := strings.TrimSuffix(strings.TrimPrefix(git, `https://gitlab.com`), ".git")
		raw = project + "/raw/"
		base := project + "/raw/" + patch + "/"
		db := dao.NewLinkDB("db")
		defer db.Close()
		fmt.Println(base)
		iter := db.NewIterator(util.BytesPrefix([]byte(base)), nil)
		ch := make(chan int, 100)
		for iter.Next() {
			key := string(iter.Key())
			value := string(iter.Value())
			if value == "1" && !db.GetBool("check"+key) {
				ch <- 1
				wg.Add(1)
				go vertify(key, client, db, ch)
			}
		}
		iter.Release()
		if iter.Error() != nil {
			fmt.Println("error")
		}
		wg.Wait()
	},
}

func vertify(url string, client http.Client, db *dao.LinkDB, ch chan int) {
	path := strings.TrimPrefix(url, raw)
	res, err := client.Get("https://gitlab.com/TeeFirefly/FireNow-Nougat/raw/" + path)

	if err != nil {
		fmt.Println(err)
		<-ch
		wg.Done()
		return
	}
	defer res.Body.Close()

	info, err := os.Stat(path)
	if err != nil {
		fmt.Println(err)
		<-ch
		wg.Done()
		return
	}
	if info.Size() == res.ContentLength {
		db.PutBool("check"+url, true)
	} else {
		fmt.Println(path, info.Size(), res.ContentLength, info.Size() == res.ContentLength)
		db.PutBool(url, false)
	}
	<-ch
	wg.Done()
}
