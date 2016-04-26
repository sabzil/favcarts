package favcarts

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"appengine"
	"appengine/datastore"
)

type UserItemPage struct {
	KeyId      int64
	ContentsId int
	Idx        int
	Title      string
	Chk        bool
	ReadCnt    int
	UserId     string
}

type UserTocItemPage struct {
	Ctnts    UserContentsIdx
	ItemPage []UserItemPage
}

func userTocItemHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mySession")
	email := session.Values["sid"]

	if email == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
	} else {
		contentsId := r.FormValue("id")
		nId, _ := strconv.Atoi(contentsId)

		c := appengine.NewContext(r)

		q_userContentsIdx := datastore.NewQuery("UserContentsIdx").Filter("Id =", nId)
		var cnts []UserContentsIdx
		_, userContentsIdx_err := q_userContentsIdx.GetAll(c, &cnts)
		if userContentsIdx_err != nil {
			http.Error(w, userContentsIdx_err.Error(), http.StatusInternalServerError)
			return
		}

		q_userTocItem := datastore.NewQuery("UserTocItem").Filter("ContentsId =", nId).Filter("UserId =", email.(string)).Order("Idx")
		var userTocItem []*UserTocItem
		k, userTocItem_err := q_userTocItem.GetAll(c, &userTocItem)
		if userTocItem_err != nil {
			http.Error(w, userTocItem_err.Error(), http.StatusInternalServerError)
			return
		}

		itms := make([]UserItemPage, len(k))

		for key, val := range k {
			itms[key].KeyId = val.IntID()
			itms[key].ContentsId = userTocItem[key].ContentsId
			itms[key].Idx = userTocItem[key].Idx
			itms[key].Title = userTocItem[key].Title
			itms[key].Chk = userTocItem[key].Chk
			itms[key].ReadCnt = userTocItem[key].ReadCnt
			itms[key].UserId = userTocItem[key].UserId

		}

		tocItemPage := &UserTocItemPage{Ctnts: cnts[0], ItemPage: itms}

		userTocItemPage := template.Must(template.New("userTocItem").Parse(userTocItemTemplate))
		userTocItemPage.Execute(w, tocItemPage)
	}
}

func tocUpdateHandler(w http.ResponseWriter, r *http.Request) {

}

func tocManipulatorHandler(w http.ResponseWriter, r *http.Request) {
	var redirectUrl string
	session, _ := store.Get(r, "mySession")
	if session.Values["sid"] == nil {
		redirectUrl = "/signin/"
	} else {
		contentsid, _ := strconv.Atoi(r.FormValue("contentsid"))
		idx, _ := strconv.Atoi(r.FormValue("idx"))
		chk, _ := strconv.ParseBool(r.FormValue("chk"))
		keyid, str_err := strconv.ParseInt(r.FormValue("keyid"), 10, 64)
		if str_err != nil {
			http.Error(w, str_err.Error(), http.StatusInternalServerError)
			return
		}
		var complete bool
		if chk == true {
			complete = false
		} else if chk == false {
			complete = true
		}
		readcnt, _ := strconv.Atoi(r.FormValue("readcnt"))

		c := appengine.NewContext(r)

		var item UserTocItem
		item.ContentsId = contentsid
		item.Idx = idx
		item.Title = r.FormValue("title")
		item.Chk = complete
		item.ReadCnt = readcnt
		item.UserId = r.FormValue("userid")
		k := userTocItemTableKey(c, keyid)

		_, err := datastore.Put(c, k, &item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		redirectUrl = fmt.Sprintf("/rating/?id=%d&idx=%d", contentsid, idx)

	}

	defer http.Redirect(w, r, redirectUrl, http.StatusFound)
}

const userTocItemTemplate = `
<!DOCTYPE html>
<html lang="ko">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">

	<title>favcarts</title>
	<link href="../static/css/bootstrap.min.css" rel="stylesheet">
</head>
<body>
	<div class="navbar navbar-inverse navbar-fixed-top" role="navigation">
		<div class="container">
			<div class="navbar-header">
				<button type="button" class="navbar-toggle" data-toggle="collapse" data-target=".navbar-collapse">
					<span class="sr-only">Toggle navigation</span>
					<span class="icon-bar"></span>
					<span class="icon-bar"></span>
					<span class="icon-bar"></span>
				</button>
				<a class="navbar-brand" href="/">FavCarts</a>
			</div>
			<div class="collapse navbar-collapse">
				<ul class="nav navbar-nav navbar-right">
					<li><a href="/">Home</a></li>
					<li><a href="/signout/">signout</a></li>
					<li><a href="/myinfo/">my</a></li>
					<li><a href="/admin/">admin</a></li>
				</ul>
			</div>
		</div>
	</div>

	<div class="container">
		<br>
		<div class="page-header">
			<div class="row">
			<div class="col-md-10"><a href="/userContentList">내 카트</a>&raquo;<h1>{{.Ctnts.Title}}</h1></div>
			<div class="col-md-2">
				<form action="." method="post">
					<input type="hidden" name="id" value="{{.Ctnts.Id}}">
					<input type="hidden" name="ver" value="{{.Ctnts.Version}}">
					<button class="btn btn-xs btn-default btn-default" type="submit">update</button>
				</form>
			</div>
			</div>
		</div>
		<div class="row">
				{{with .ItemPage}}
				{{range .}}
					<div class="col-md-11">{{if .Chk}}<strike>{{.Title}}</strike>{{else}}{{.Title}}{{end}}</div>
					<div class="col-md-1">					
						<form action="/tocManipulator/" method="post">
						<input type="hidden" name="keyid" value={{.KeyId}}>
						<input type="hidden" name="title" value="{{.Title}}">
						<input type="hidden" name="chk" value="{{.Chk}}">
						<input type="hidden" name="readcnt" value="{{.ReadCnt}}">
						<input type="hidden" name="idx" value="{{.Idx}}">
						<input type="hidden" name="userid" value="{{.UserId}}">
						<input type="hidden" name="contentsid" value="{{.ContentsId}}">
						{{if .Chk}}
							<button class="btn btn-xs btn-default btn-block" type="submit">√</button>
						{{else}}
							<button class="btn btn-xs btn-default btn-block" type="submit">□</button>
						{{end}}
						</form>
					</div>
				{{end}}
				{{end}}
		</div>
		</form>
	</div>
</body>
</html>`
