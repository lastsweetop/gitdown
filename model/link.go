package model

import "github.com/lastsweetop/gitdown/dao"

type LinkP struct {
	Path string
	Db   *dao.LinkDB
}
