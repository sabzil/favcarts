package favcarts

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"

	"appengine"
	"appengine/datastore"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		signPage := template.Must(template.New("signin").Parse(signinTemplate))
		signPage.Execute(w, nil)
	} else if r.Method == "POST" {
		userSession, _ := store.Get(r, "mySession")
		email := r.FormValue("email")
		pw := r.FormValue("password")

		if len(email) == 0 || len(pw) == 0 {
			fmt.Fprint(w, "<script>alert('email과 password를 입력해주세요.');location.href='/signin/';</script>")
			return
		}

		c := appengine.NewContext(r)
		q := datastore.NewQuery("Users").Filter("Email=", email)
		var data []User
		_, err := q.GetAll(c, &data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(data) == 0 {
			fmt.Fprint(w, "<script>alert('email이 존재하지 않습니다.');location.href='/signin/';</script>")
			return

		}

		hasher := sha1.New()
		hasher.Write([]byte(pw))
		hashData := hasher.Sum(nil)
		sha := base64.URLEncoding.EncodeToString(hashData)

		if data[0].Password == sha {
			userSession.Values["sid"] = email
			userSession.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			fmt.Fprint(w, "<script>alert('password를 확인해 주세요');location.href='/signin/';</script>")
		}

	}
}

const signinTemplate = `

<!DOCTYPE html>
<html lang="ko">
<head>
	<meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
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
					<li><a href="/signin/">signin</a></li>
					<li><a href="/admin/">admin</a></li>
				</ul>
			</div>
		</div>
	</div>

    <div class="container">
		<div class="page-header">
			<br><h1>로그인</h1>
		</div>

		<div class="container">
		<form class="form-signin" role="form" action="/signin/" method="post">
			<!-- <h2 class="form-signin-heading">Please sign in</h2> -->
			<input type="email" class="form-control" placeholder="Email address" name="email" required autofocus>
			<input type="password" class="form-control" placeholder="Password" name="password" required>			
			<button class="btn btn-lg btn-primary btn-block" type="submit">로그인</button>
		</form>		
		<a href="/signup/" class="btn btn-defalut btn-sm pull-right" role="button">회원가입</a>
		</div>
	</div>
</body>
</html>
`
