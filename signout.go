package favcarts

import (
	"log"
	"net/http"
)

func signoutHandler(w http.ResponseWriter, r *http.Request) {
	userSession, err := store.Get(r, "mySession")
	if err != nil {
		log.Fatal(err)
	}

	for key, _ := range userSession.Values {
		delete(userSession.Values, key)
	}

	userSession.Save(r, w)

	http.Redirect(w, r, "/signin/", http.StatusFound)
}
