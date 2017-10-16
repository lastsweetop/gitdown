package model

import "gitdown/dao"

type LinkP struct {
	Path string
	Db   *dao.LinkDB
}
