package controller

import (
	"fmt"
	"html/template"
	"model"
	"net/http"
	"strconv"
	"strings"
	"utils"
)

var cachedTags map[string]model.Tag = nil

type Result struct {
	Blogs map[string]model.Blog
	Next  int
	Pre   int
	Admin bool
	Tags  map[string]model.Tag
}

func unescaped(x string) interface{} { return template.HTML(x) }
func haspre(pre int) bool {
	return pre != 0
}

func pageing(page string) (n, p, i int) {
	if page == "" {
		i = 1
	} else {
		i, _ = strconv.Atoi(page)
	}
	n = i + 1
	p = i - 1
	if i >= 1 {
		i = (i - 1) * utils.Limit
	}

	return n, p, i
}

func List(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	page := r.Form.Get("page")

	n, p, i := pageing(page)

	db, _ := model.Init()

	blogs := db.List(i, utils.Limit)

	if cachedTags == nil {
		cachedTags = db.FindTag()
	}

	cookie, _ := r.Cookie("blog")
	results := new(Result)
	if cookie != nil && cookie.Value == "anything" {
		results.Admin = true
	} else {
		results.Admin = false
	}

	results.Blogs = blogs
	results.Next = n
	results.Pre = p
	results.Tags = cachedTags

	t := template.New("list.html")
	t = t.Funcs(template.FuncMap{"unescaped": unescaped, "haspre": haspre})
	t, err := t.ParseFiles(utils.View_addr + "/list.html")

	// t, err := template.ParseFiles(utils.View_addr + "/list.html")
	err = t.Execute(w, &results)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}

func Blog(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("blog")
	if cookie != nil {
		if cookie.Value == "anything" {
			t, err := template.ParseFiles(utils.View_addr + "/blogpublish.html")
			err = t.Execute(w, nil)
			if err != nil {
				fmt.Println("Fatal error ", err.Error())
			}
		} else {
			http.Redirect(w, r, "/index", http.StatusForbidden)
		}
	} else {
		http.Redirect(w, r, "/index", http.StatusForbidden)
	}

}

func Save(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	blog := new(model.Blog)

	blog.Title = r.Form["blogTitle"][0]
	blog.Content = r.Form["content"][0]
	db, _ := model.Init()

	blogid := db.Add(blog)

	tagName := r.Form["tagName"][0]
	tagNames := strings.Split(tagName, ",")

	db.AddTag(blogid, tagNames)
	db.AddCtr(blogid) // 初始化点击率

	cachedTags = nil

	http.Redirect(w, r, "/index", http.StatusFound)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form.Get("id")

	db, _ := model.Init()
	idint, _ := strconv.Atoi(id)
	blog := db.FindById(idint)
	blog = db.PutTag(blog)
	t, err := template.ParseFiles(utils.View_addr + "/edit.html")
	t.Execute(w, blog)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}

func DoEdit(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	blog := new(model.Blog)

	blog.Title = r.Form["blogTitle"][0]
	blog.Content = r.Form["content"][0]
	id := r.Form["id"][0]
	blog.Id, _ = strconv.Atoi(id)
	db, _ := model.Init()
	db.Edit(blog)

	tagName := r.Form["tagName"][0]
	// fmt.Println(tagName)
	tagNames := strings.Split(tagName, ",")
	// fmt.Println(len(tagNames))
	db.EditTag(blog.Id, tagNames)

	cachedTags = nil

	http.Redirect(w, r, "/index", http.StatusFound)

}

func Delete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form.Get("id")
	fmt.Println("id %s", id)
	idint, _ := strconv.Atoi(id)
	db, _ := model.Init()
	db.Delete(idint)

	db.DelTagByBlogId(idint)

	cachedTags = nil

	http.Redirect(w, r, "/index", http.StatusFound)
}

func FindByTag(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tagname := r.Form.Get("tagName")
	page := r.Form.Get("page")
	n, p, i := pageing(page)

	db, _ := model.Init()
	blogs := db.FindByTag(i, utils.Limit, tagname)

	cookie, _ := r.Cookie("blog")
	results := new(Result)
	if cookie != nil && cookie.Value == "anything" {
		results.Admin = true
	} else {
		results.Admin = false
	}
	results.Blogs = blogs
	results.Next = n
	results.Pre = p
	results.Tags = cachedTags

	t := template.New("list.html")
	t = t.Funcs(template.FuncMap{"unescaped": unescaped, "haspre": haspre})
	t, err := t.ParseFiles(utils.View_addr + "/list.html")
	err = t.Execute(w, &results)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}

func Detail(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form.Get("id")

	db, _ := model.Init()
	idint, _ := strconv.Atoi(id)
	blog := db.FindById(idint)
	blog = db.PutCount(blog)

	go db.Click(idint) // 记录一次点击

	t := template.New("detail.html")
	t = t.Funcs(template.FuncMap{"unescaped": unescaped})
	t, err := t.ParseFiles(utils.View_addr + "/detail.html")
	err = t.Execute(w, blog)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
	}
}
