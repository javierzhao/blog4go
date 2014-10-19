package controller

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"model"
	"net/http"
	"time"
	"utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(utils.View_addr + "/login.html")
	err = t.Execute(w, nil)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}

func DoLogin(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	username := r.Form["username"][0]
	password := r.Form["password"][0]

	h := md5.New()
	h.Write([]byte(password))

	pass_word := hex.EncodeToString(h.Sum(nil))

	db, _ := model.Init()
	isUser := db.IsUser(username, pass_word)

	if isUser {
		expiration := time.Now()
		expiration = expiration.AddDate(1, 0, 0)
		cookie := http.Cookie{Name: "blog", Value: "anything", Expires: expiration}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/blog", http.StatusFound)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}

}
