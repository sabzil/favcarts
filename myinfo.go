package favcarts

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"appengine"
	"appengine/datastore"
)

type Myinfo struct {
	KeyId    int64
	Name     string
	Email    string
	Password string
}

func myInfoHandler(w http.ResponseWriter, r *http.Request) {
	userSession, _ := store.Get(r, "mySession")
	email := userSession.Values["sid"]

	if r.Method == "GET" {
		if email == nil {
			http.Redirect(w, r, "/signin/", http.StatusFound)
		} else {
			c := appengine.NewContext(r)

			userQuery := datastore.NewQuery("Users").Filter("Email=", email)

			var user []User
			userKey, userErr := userQuery.GetAll(c, &user)
			if userErr != nil {
				http.Error(w, userErr.Error(), http.StatusInternalServerError)
				return
			}

			my := &Myinfo{KeyId: userKey[0].IntID(), Name: user[0].Name, Email: user[0].Email, Password: user[0].Password}

			myInfoPage := template.Must(template.New("myInfo").Parse(myInfoTemplate))
			myInfoPage.Execute(w, my)
		}

	} else if r.Method == "POST" {
		if email == nil {
			http.Redirect(w, r, "/signin/", http.StatusFound)
		} else {
			emailId := r.FormValue("email")
			nameId := r.FormValue("name")
			password := r.FormValue("password1")
			confirm := r.FormValue("password2")
			keyid, err := strconv.ParseInt(r.FormValue("keyid"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("key : %d", keyid)
			log.Println(emailId)
			log.Println(nameId)
			log.Printf("password : %s", password)

			if len(password) == 0 || len(confirm) == 0 {
				fmt.Fprint(w, "<script>alert('password를 확인해 주세요(0)');location.href='/myinfo/';</script>")
			}

			if password != confirm {
				fmt.Fprint(w, "<script>alert('password를 확인해 주세요(1)');location.href='/myinfo/';</script>")
			}

			c := appengine.NewContext(r)
			k := userTableKey(c, keyid)

			hasher := sha1.New()
			hasher.Write([]byte(password))
			hashData := hasher.Sum(nil)
			sha := base64.URLEncoding.EncodeToString(hashData)

			user := &User{
				Name:     nameId,
				Email:    emailId,
				Password: sha}

			_, err = datastore.Put(c, k, user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/myinfo/", http.StatusFound)

		}
	}
}

func deleteAccountHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		userSession, _ := store.Get(r, "mySession")
		email := userSession.Values["sid"]

		if email == nil {
			http.Redirect(w, r, "/signin/", http.StatusFound)
		} else {
			keyid, err := strconv.ParseInt(r.FormValue("keyid"), 10, 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			emailid := r.FormValue("email")
			log.Printf("keyid : %d, email:%s", keyid, emailid)

			c := appengine.NewContext(r)
			k := userTableKey(c, keyid)

			err = datastore.Delete(c, k)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			////////////usercontentslist와 usertocitem을 삭제
			q_userTocItem := datastore.NewQuery("UserTocItem").Filter("UserId =", email.(string)).KeysOnly()
			userTocItemKey, userTocItemErr := q_userTocItem.GetAll(c, nil)
			if userTocItemErr != nil {
				http.Error(w, userTocItemErr.Error(), http.StatusInternalServerError)
				return
			}
			for _, itmKey := range userTocItemKey {
				datastore.Delete(c, userTocItemTableKey(c, itmKey.IntID()))
			}

			q_userContentsIdx := datastore.NewQuery("UserContentsIdx").Filter("UserId =", email.(string)).KeysOnly()
			userContentsIdxKey, userContentsIdxErr := q_userContentsIdx.GetAll(c, nil)
			if userContentsIdxErr != nil {
				http.Error(w, userTocItemErr.Error(), http.StatusInternalServerError)
				return
			}
			for _, contentsKey := range userContentsIdxKey {
				datastore.Delete(c, userContentsIdxTableKey(c, contentsKey.IntID()))
			}

			http.Redirect(w, r, "/signout/", http.StatusFound)

		}

	}

}

const myInfoTemplate = `
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
					<li><a href="/signout/">signout</a></li>
					<li class="active"><a href="/myinfo/">my</a></li>
					<li><a href="/admin/">admin</a></li>
				</ul>
			</div>
		</div>
	</div>

	<div class="container">
		<div class="page-header">
			<br><h1>회원정보</h1>
		</div>

		<div class="container">
			<form class="form-horizontal" action="/myinfo/" method="post">
				<fieldset>
					<legend><b>비밀번호 변경</b></legend>
					<div class="control-group">
						<div class="controls">
						<input type="hidden" id="keyid" name="keyid" value="{{.KeyId}}">
						</div>
					</div>
					<div class="control-group">
						<label class="control-label" for="email">email</label>
						<div class="controls">
							<input type="hidden" id="email" name="email" value="{{.Email}}">
							<input type="text" class="input-medium disabled" id="email" value="{{.Email}}" disabled>
						</div>
					</div>
					<div class ="control-group">
						<label class="control-label" for="name">name</label>
						<div class="controls">
							<input type="hidden" id="name" name="name" value="{{.Name}}">
							<input type="text" class="input-medium disable" id="name" value="{{.Name}}" disabled>
						</div>
					</div>
					<div class="control-group">
						<label class="control-label" for="password1">password</label>
						<div class="controls">
							<input type="password" class="input-medium" id="passwordi1" name="password1" value="">
						</div>
					</div>
					<div class="control-group">
						<label class="control-label" for="password2">confirm</label>
						<div class="control-group">
							<input type="password" class="input-medium" id="password2" name="password2" value="">
						</div>
					</div>
					<br/>
					<div class="form-actions">
						<button type="submit" class="btn btn-primary">변경</button>
						<a href="/" class="btn">취소</a>
					</div>
				</fieldset>
			</form>
		</div>
		<br/>
		<div class="container">
			<form class="form-horizontal" action="/deleteAccount/" method="post">
	
				<fieldset>
					<legend><b>계정삭제</b></legend>
					<div class="control-group">
						<div class="controls">
							<input type="hidden" id="keyid" name="keyid" value="{{.KeyId}}">
						</div>
					</div>
					<div class="control-group">
						<div class="controls">
							<input type="hidden" id="email" name="email" value="{{.Email}}">
						</div>
					</div>
					<div class="form-actions">
						<button type="submit" class="btn btn-primary">탈퇴</button>
					</div>
				</fieldset>
			</form>
		</div>
			
	</div>

</body>
</html>

`
