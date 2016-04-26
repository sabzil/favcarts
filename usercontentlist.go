package favcarts

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"
)

func userContentListHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mySession")
	email := session.Values["sid"]

	if email == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
	} else {
		c := appengine.NewContext(r)
		q_userContentsidx := datastore.NewQuery("UserContentsIdx").Filter("UserId =", email.(string))
		var userCnts []UserContentsIdx
		_, userCnts_err := q_userContentsidx.GetAll(c, &userCnts)
		if userCnts_err != nil {
			http.Error(w, userCnts_err.Error(), http.StatusInternalServerError)
			return
		}

		userContentListPage := template.Must(template.New("userContentList").Parse(userContentListTemplate))
		userContentListPage.Execute(w, userCnts)
	}
}

const userContentListTemplate = `
<!DOCTYPE html>
<html lang="ko">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">

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
			<br>
			<h1>내 카트</h1>
		</div>
		<div id="row">
			<ul>
				{{range .}}
				<li><a href="/userTocItem/?id={{.Id}}">{{.Title}}</a></li>
				{{end}}
			</ul>
		</div>
	</div>
	
</body>
</html>
`
