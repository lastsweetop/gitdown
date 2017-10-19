package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/lastsweetop/gitdown/dao"
	"github.com/spf13/cobra"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func init() {
	vertifyCmd.Flags().StringVarP(&git, "git", "g", "https://gitlab.com/TeeFirefly/FireNow-Nougat/", "git url")
	vertifyCmd.Flags().StringVarP(&patch, "patch", "p", "master", "patch")
	RootCmd.AddCommand(vertifyCmd)
}

var vertifyCmd = &cobra.Command{
	Use:   "vertify",
	Short: "vertify the repo",
	Run: func(cmd *cobra.Command, args []string) {
		project := strings.TrimSuffix(strings.TrimPrefix(git, `https://gitlab.com`), ".git")
		raw = project + "/raw/"
		base := project + "/raw/" + patch + "/"
		db := dao.NewLinkDB("db")
		defer db.Close()
		iter := db.NewIterator(util.BytesPrefix([]byte(base)), nil)
		ch := make(chan int, 200)
		for iter.Next() {
			key := string(iter.Key())
			value := string(iter.Value())
			if value == "1" {
				ch <- 1
				wg.Add(1)
				go vertify(key, ch)
			}
		}
		iter.Release()
		if iter.Error() != nil {
			fmt.Println("error")
		}
		wg.Wait()
	},
}

func vertify(url string, ch chan int) {
	path := strings.TrimPrefix(url, raw)
	if _, err := os.Stat(path); err != nil {
		fmt.Println(path)
	}
	<-ch
	wg.Done()
}
