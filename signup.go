package favcarts

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"text/template"

	"appengine"
	"appengine/datastore"
)

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		signupPage := template.Must(template.New("signup").Parse(signupTemplate))
		signupPage.Execute(w, nil)
	} else if r.Method == "POST" {
		name := r.FormValue("name")
		id := r.FormValue("email")
		pw := r.FormValue("password")
		if len(name) == 0 || len(id) == 0 || len(pw) == 0 {
			fmt.Fprint(w, "<script>alert('name, email 또는  password를 입력해주세요.');location.href='/signup/';</script>")
			return
		}

		c := appengine.NewContext(r)

		validQuery := datastore.NewQuery("Users").Filter("Email =", id).KeysOnly()
		validKeys, validErr := validQuery.GetAll(c, nil)
		if validErr != nil {
			http.Error(w, validErr.Error(), http.StatusInternalServerError)
			return
		}
		if len(validKeys) > 0 {
			fmt.Fprint(w, "<script>alert('동일한 email로 가입이 되어 있습니다.');location.href='/signup/';</script>")
			return
		}

		hasher := sha1.New()
		hasher.Write([]byte(pw))
		hashData := hasher.Sum(nil)
		sha := base64.URLEncoding.EncodeToString(hashData)

		user := &User{
			Name:     name,
			Email:    id,
			Password: sha,
		}

		key := datastore.NewIncompleteKey(c, "Users", nil)
		_, err := datastore.Put(c, key, user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/signin/", http.StatusFound)
	}
}

const signupTemplate = `
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
				<br><h1>회원가입</h1>
			</div>
			<div class="container">
			<form class="form-signup" role="form" action="/signup/" method="post">
				<!-- <h2 class="form-signup-heading">Please sign up</h2> -->
				<input type="name" class="form-control" placeholder="Name" name="name" required autofocus>
				<input type="email" class="form-control" placeholder="Email address" name="email" required>
				<input type="password" class="form-control" placeholder="Password" name="password" required>			
				<button class="btn btn-lg btn-primary btn-block" type="submit">회원가입</button>
			</form>
			</div>
		</div>
		
	</body>
</html>
`
