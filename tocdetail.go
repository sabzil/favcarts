package favcarts

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"

	"strconv"
)

type DetailTocItem struct {
	Item     TocItem
	Keywords []string
}

type Detail struct {
	Title       string
	ContentsIdx int
	Desc        string
	Items       []DetailTocItem
}

func tocdetailHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mySession")
	if session.Values["sid"] == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
	} else {
		c := appengine.NewContext(r)

		var detail Detail
		contentsId := r.FormValue("id")
		nId, _ := strconv.Atoi(contentsId)

		var ctntIdx []ContentsIdx
		q_title := datastore.NewQuery("ContentsIdx").Filter("Id =", nId)
		_, title_err := q_title.GetAll(c, &ctntIdx)
		if title_err != nil {
			http.Error(w, title_err.Error(), http.StatusInternalServerError)
			return
		}

		detail.Title = ctntIdx[0].Title
		detail.ContentsIdx = ctntIdx[0].Id

		q := datastore.NewQuery("TocItem").Filter("ContentsId =", nId).Order("Idx")

		var toc []TocItem
		_, err := q.GetAll(c, &toc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//log.Printf("len:%d", len(toc))
		detailItem := make([]DetailTocItem, len(toc))
		for idx, itm := range toc {
			detailItem[idx].Item = itm

			tagQuery := datastore.NewQuery("Tag").Filter("ContentsId =", nId).Filter("ItemIdx =", itm.Idx)
			var tag []Tag
			_, err := tagQuery.GetAll(c, &tag)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if len(tag) > 0 {
				detailItem[idx].Keywords = make([]string, len(tag))
				for tagIdx, tagItm := range tag {
					detailItem[idx].Keywords[tagIdx] = tagItm.Keyword
				}
			}

		}
		detail.Items = detailItem

		tocdetailPage := template.Must(template.New("tocdetail").Parse(tocdetailTemplate))
		tocdetailPage.Execute(w, detail)
	}
}

const tocdetailTemplate = `
<!DOCTYPE html>
<html lang="ko">
	<head>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
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

			<div class="page-header">
				<div class="row">
					<br>
					<div class="col-md-10">
						<h1>{{.Title}}</h1>
					</div>
					<div class="col-md-2">
						<form action="/addMyToc/" method="post">
							<input type="hidden" name="contentsIdx" value="{{.ContentsIdx}}">
							<button class="btn btn-default" type="submit">카트에 담기</button>
						</form>
					</div>
				</div>
			</div>
			<div class="row">
				<ul>
					{{with .Items}}
					{{range .}}
						<li>{{.Item.Title}}<br/>
						<p>
						{{range .Keywords}}
						{{.}} 
						{{end}}</p>
						</li>
					{{end}}
					{{end}}
				</ul>
			</div>
		</div>
	</body>
</html>
`

func addMyTocHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		session, _ := store.Get(r, "mySession")
		email := session.Values["sid"]
		if email == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
		} else {
			contentsIdx := r.FormValue("contentsIdx")
			//ContetnsIdx테이블에서  idx를 갖고 온다.
			c := appengine.NewContext(r)
			nId, _ := strconv.Atoi(contentsIdx)

			//user테이블에 같은 idx로 추가된게 있는지 확인한다.
			queryCheck := datastore.NewQuery("UserContentsIdx").Filter("UserId=", email).Filter("Id=", nId).KeysOnly()
			chkKeys, chkErr := queryCheck.GetAll(c, nil)
			if chkErr != nil {
				http.Error(w, chkErr.Error(), http.StatusInternalServerError)
				return
			}
			if len(chkKeys) > 0 {
				http.Redirect(w, r, "/userContentList/", http.StatusFound)
			} else {

				q_ContentsIdx := datastore.NewQuery("ContentsIdx").Filter("Id =", nId)
				var cnts []ContentsIdx
				_, contentsIdx_err := q_ContentsIdx.GetAll(c, &cnts)
				if contentsIdx_err != nil {
					http.Error(w, contentsIdx_err.Error(), http.StatusInternalServerError)
					return
				}

				q_TocItem := datastore.NewQuery("TocItem").Filter("ContentsId =", nId)
				var tocItems []TocItem
				_, tocItems_err := q_TocItem.GetAll(c, &tocItems)
				if tocItems_err != nil {
					http.Error(w, tocItems_err.Error(), http.StatusInternalServerError)
					return
				}

				userContentsIdx := &UserContentsIdx{
					CategoryId:  cnts[0].CategoryId,
					Id:          cnts[0].Id,
					Title:       cnts[0].Title,
					Description: cnts[0].Description,
					Author:      cnts[0].Author,
					Publisher:   cnts[0].Publisher,
					UserId:      email.(string),
					Version:     cnts[0].Version}

				userContentsIdx_key := datastore.NewIncompleteKey(c, "UserContentsIdx", nil)
				_, userContentsIdx_err := datastore.Put(c, userContentsIdx_key, userContentsIdx)
				if userContentsIdx_err != nil {
					http.Error(w, userContentsIdx_err.Error(), http.StatusInternalServerError)
					return
				}

				userTocItem_key := datastore.NewIncompleteKey(c, "UserTocItem", nil)
				for _, itm := range tocItems {
					userTocItems := &UserTocItem{
						ContentsId: itm.ContentsId,
						Idx:        itm.Idx,
						Title:      itm.Title,
						Chk:        false,
						ReadCnt:    0,
						UserId:     email.(string)}

					_, userTocItem_err := datastore.Put(c, userTocItem_key, userTocItems)
					if userTocItem_err != nil {
						http.Error(w, userTocItem_err.Error(), http.StatusInternalServerError)
						return
					}
				}

			}

			http.Redirect(w, r, "/userContentList/", http.StatusFound)
		}
	}
}
