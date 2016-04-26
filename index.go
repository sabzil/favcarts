package favcarts

import (
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"

	"appengine"
	"appengine/datastore"
)

var store = sessions.NewCookieStore([]byte("blah, blah(알아서 추가)"))

func init() {
	store.Options = &sessions.Options{
		Domain:   "localhosti:8080",
		Path:     "/",
		MaxAge:   3600 * 1, //1 hour
		HttpOnly: true,
	}
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/signin/", signinHandler)
	http.HandleFunc("/signup/", signupHandler)
	http.HandleFunc("/signout/", signoutHandler)
	http.HandleFunc("/deleteAccount/", deleteAccountHandler)

	http.HandleFunc("/admin/", adminHandler)
	http.HandleFunc("/adminCategory/", adminCategoryHandler)
	http.HandleFunc("/adminContentsIdx/", adminContentsIdxHandler)

	http.HandleFunc("/adminOnePiece1/", adminOnePiece1Handler)
	http.HandleFunc("/adminOnePiece2/", adminOnePiece2Handler)
	http.HandleFunc("/adminOnePiece3/", adminOnePiece3Handler)
	http.HandleFunc("/adminOnePiece4/", adminOnePiece4Handler)
	http.HandleFunc("/adminOnePiece5/", adminOnePiece5Handler)
	http.HandleFunc("/adminOnePiece6/", adminOnePiece6Handler)
	http.HandleFunc("/adminOnePiece7/", adminOnePiece7Handler)

	http.HandleFunc("/adminTocItem1/", adminTocItem1Handler)
	http.HandleFunc("/adminTocItem2/", adminTocItem2Handler)
	http.HandleFunc("/adminTocItem3/", adminTocItem3Handler)

	http.HandleFunc("/tocdetail/", tocdetailHandler)
	http.HandleFunc("/addMyToc/", addMyTocHandler)
	http.HandleFunc("/userContentList/", userContentListHandler)
	http.HandleFunc("/userTocItem/", userTocItemHandler)
	http.HandleFunc("/tocUpdateManipulator/", tocUpdateManipulatorHandler)
	http.HandleFunc("/tocManipulator/", tocManipulatorHandler)
	http.HandleFunc("/rating/", ratingHandler)
	http.HandleFunc("/ratingManipulator/", ratingManipulatorHandler)
	http.HandleFunc("/myinfo/", myInfoHandler)
}

var indexTemplate = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="ko">
<head>
	<meta charset="utf-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>favcarts</title>
	<link href="static/css/bootstrap.min.css" rel="stylesheet">
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
					<li class="active"><a href="/">Home</a></li>
					{{if .Signin}}
						<li><a href="/signout/">signout</a></li>
					{{else}}
						<li>a href="/singin/">singin</a></li>
					{{end}}
					<li><a href="/myinfo/">my</a></li>
					<li><a href="/admin/">admin</a></li>
				</ul>
			</div>
		</div>
	</div>

	<div class="container">
		<div class="page-header">
			<br>
			<h1>Table of Contetns</h1>
		</div>
		<div class="row">
			<div class="col-md-8">
				<div class="row">
					{{with .ContentsList}}
					{{range .}}
					<div class="col-md-4">
						<h2>{{.Title}}</h2>
						<p>{{.Description}}</p>
						<p><a class="btn btn-default" href="/tocdetail/?id={{.Id}}">자세히 &raquo;</a></p>
					</div>
					{{end}}
					{{end}}
				</div>
			</div>
			<div class="col-md-4">
				<h4><a href="/userContentList/">내 카트</a></h4>
				<ul>
					{{with .UserContentsList}}
					{{range .}}
						<li><a href="/userTocItem/?id={{.Id}}">{{.Title}}</a></li>
					{{end}}
					{{end}}
				</ul>
			</div>
		</div>
	</div>

</body>
</html>

`))

type IndexPage struct {
	Signin           bool
	ContentsList     []ContentsIdx
	UserContentsList []UserContentsIdx
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mySession")
	email := session.Values["sid"]
	if email == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
	} else {
		c := appengine.NewContext(r)
		q := datastore.NewQuery("ContentsIdx")
		var data []ContentsIdx
		_, err := q.GetAll(c, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user_q := datastore.NewQuery("UserContentsIdx").Filter("UserId=", email)
		var userData []UserContentsIdx
		_, err = user_q.GetAll(c, &userData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		idxPage := &IndexPage{Signin: true, ContentsList: data, UserContentsList: userData}

		indexTemplate.Execute(w, idxPage)
	}
}
