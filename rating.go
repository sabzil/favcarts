package favcarts

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

type Rating struct {
	Signin bool
	Id     int
	Idx    int
	Rate   int
}

func ratingHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "mySession")
	email := session.Values["sid"]

	if email == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
	} else {
		id, _ := strconv.Atoi(r.FormValue("id"))
		idx, _ := strconv.Atoi(r.FormValue("idx"))

		rateInfo := &Rating{Id: id, Idx: idx}
		ratingPage := template.Must(template.New("rating").Parse(ratingTemplate))
		ratingPage.Execute(w, rateInfo)
	}

}

const ratingTemplate = `
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
		<div class ="page-header">
			<br><h1>당신의 점수는 몇 점?</h1>
		</div>
		<div class="container">
			<form action="/ratingManipulator/" method="post">
				<input type="hidden" name="id" value="{{.Id}}">
				<input type="hidden" name="idx" value="{{.Idx}}">
				<select class="form-control" name="rate">
					<option value="1">1</option>
					<option value="2">2</option>
					<option value="3">3</option>
					<option value="4">4</option>
					<option value="5">5</option>
				</select>
				<button class="btn btn-xs btn-default pull-right" type="submit">ok</button>
				<!-- <button class="btn btn-xs btn-default pull-right" type="button" href="/usertocitem/?id={{.Id}}">cancel</button> -->
				<a href="/userTocItem/?id={{.Id}}" class="btn btn-xs btn-default pull-right">cancel</a>
			</form>
		</div>
	</div>
</body>
</html>
`

func ratingManipulatorHandler(w http.ResponseWriter, r *http.Request) {
	contentsid, _ := strconv.Atoi(r.FormValue("id"))
	rate, _ := strconv.Atoi(r.FormValue("rate"))
	log.Printf("rate : %d", rate)

	http.Redirect(w, r, fmt.Sprintf("/userTocItem/?id=%d", contentsid), http.StatusFound)

}
