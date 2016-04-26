package favcarts

import (
	"appengine"
	"appengine/datastore"
)

func userTableKey(c appengine.Context, keyid int64) *datastore.Key {
	return datastore.NewKey(c, "Users", "", keyid, nil)
}

func contentsIdxTableKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "ContentsIdx", "contentsidx", 0, nil)
}

func categoryTableKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Category", "category", 0, nil)
}

func tocItemTableKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "TocItem", "tocitem", 0, nil)
}

func tagTableKey(c appengine.Context, keyid int64) *datastore.Key {
	return datastore.NewKey(c, "Tag", "", keyid, nil)
}

func userContentsIdxTableKey(c appengine.Context, keyid int64) *datastore.Key {
	return datastore.NewKey(c, "UserContentsIdx", "", keyid, nil)
}

func userTocItemTableKey(c appengine.Context, keyid int64) *datastore.Key {
	return datastore.NewKey(c, "UserTocItem", "", keyid, nil)
}

type User struct {
	Name     string
	Email    string
	Password string
}

type Category struct {
	Id       int
	Category string
}

type ContentsIdx struct {
	Id          int
	CategoryId  int
	Title       string
	Description string
	Author      string
	Publisher   string
	Version     string
	Score       float64
}

type TocItem struct {
	ContentsId int
	Idx        int
	Title      string
}

type Tag struct {
	Id         int
	ContentsId int
	ItemIdx    int
	Keyword    string
}

type UserContentsIdx struct {
	Id          int
	CategoryId  int
	Title       string
	Description string
	Author      string
	Publisher   string
	UserId      string
	Version     string
}

type UserTocItem struct {
	ContentsId int
	Idx        int
	Title      string
	Chk        bool
	ReadCnt    int
	UserId     string
	Memo       string
}
