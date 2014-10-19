package model

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	_ "time"
	"utils"
)

type Blog struct {
	Id       int
	Title    string
	Content  string
	Created  string
	TagNames string
	Count    string
}

type Tag struct {
	Id     int
	BlogId int
	Name   string
}

type Comment struct {
	Id      int
	BlogId  int
	Content string
	Created uint8
}

func (this *Dao) FindById(id int) *Blog {
	rows, err := this.db.Query("select id,title,content,created from blogs where id = ?", id)
	CheckErr(err)
	defer rows.Close()
	var blog = new(Blog)
	for rows.Next() {

		err = rows.Scan(&blog.Id, &blog.Title, &blog.Content, &blog.Created)
		CheckErr(err)
		break
	}

	return blog
}

func (this *Dao) PutTag(blog *Blog) *Blog {
	rows, err := this.db.Query("select name from tags where blogid = ?", blog.Id)
	CheckErr(err)
	defer rows.Close()
	tagNames := ""
	i := 1
	for ; rows.Next(); i++ {
		var name string = ""
		err = rows.Scan(&name)
		CheckErr(err)

		tagNames = tagNames + name + ","
	}

	// Go竟然没有substring的方法。。
	li := strings.LastIndex(tagNames, ",")
	l := len(tagNames)
	if li != -1 && l > 0 {
		tagNames = utils.Substr(tagNames, 0, li)
	}
	blog.TagNames = tagNames

	return blog
}

func (this *Dao) PutCount(blog *Blog) *Blog {
	rows, err := this.db.Query("select count from ctrs where blogid = ?", blog.Id)
	CheckErr(err)
	defer rows.Close()
	count := ""
	if rows.Next() {
		err = rows.Scan(&count)
		CheckErr(err)
	}

	blog.Count = count

	return blog
}

// limit = 5 ，表示每页显示多少条，page表示第几页 从0开始
func (this *Dao) List(page, limit int) map[string]Blog {
	rows, err := this.db.Query("select id,title,content,created from blogs order by created desc limit ?,?", page, limit)
	CheckErr(err)
	defer rows.Close()
	blogs := make(map[string]Blog)
	for i := 0; rows.Next(); i++ {
		var blog = new(Blog)
		err = rows.Scan(&blog.Id, &blog.Title, &blog.Content, &blog.Created)
		CheckErr(err)
		key := fmt.Sprintf("%d", i)
		blogs[key] = *blog
	}

	return blogs
}

func (this *Dao) Add(blog *Blog) int {
	// tx, _ := this.db.Begin()
	stmt, err := this.db.Prepare("INSERT blogs SET title=?,Content=?")
	CheckErr(err)
	defer stmt.Close()

	res, err1 := stmt.Exec(blog.Title, blog.Content)
	CheckErr(err1)

	lastId, _ := res.LastInsertId()
	return int(lastId)
}

func (this *Dao) Edit(blog *Blog) {
	stmt, err := this.db.Prepare("UPDATE blogs SET title=?,content=? where id = ?")
	CheckErr(err)
	defer stmt.Close()

	// fmt.Println("blog %s,%s,%s", blog.Title, blog.Content, blog.Id)

	_, err1 := stmt.Exec(blog.Title, blog.Content, blog.Id)
	CheckErr(err1)
}

func (this *Dao) Delete(id int) {
	stmt, err := this.db.Prepare("DELETE FROM blogs where id = ?")
	CheckErr(err)
	defer stmt.Close()

	_, err1 := stmt.Exec(id)
	CheckErr(err1)
}

func (this *Dao) EditTag(blogid int, name []string) {
	stmt, err := this.db.Prepare("DELETE FROM tags where blogid = ?")
	CheckErr(err)
	defer stmt.Close()

	_, err1 := stmt.Exec(blogid)
	CheckErr(err1)

	stmt, err = this.db.Prepare("INSERT tags SET blogid=?,name=?")
	CheckErr(err)

	// 实现一个变相的batch
	finish := make(chan bool)
	lenth := len(name)
	for i := 0; i < lenth; i++ {
		go func(tagName string) {
			defer func() { finish <- true }()
			if _, err = stmt.Exec(blogid, tagName); err != nil {
				fmt.Println("stmt.Exec: ", err.Error())
				return
			}
		}(name[i])
	}

	for i := 0; i < lenth; i++ {
		<-finish
	}
}

func (this *Dao) AddTag(blogid int, name []string) {
	stmt, err := this.db.Prepare("INSERT tags SET blogid=?,name=?")
	CheckErr(err)

	// 实现一个变相的batch
	finish := make(chan bool)
	lenth := len(name)
	for i := 0; i < lenth; i++ {
		go func(tagName string) {
			defer func() { finish <- true }()
			if _, err = stmt.Exec(blogid, tagName); err != nil {
				fmt.Println("stmt.Exec: ", err.Error())
				return
			}
		}(name[i])
	}

	for i := 0; i < lenth; i++ {
		<-finish
	}
}

/**
func (this *Dao) AddTag(blogid int, name string) {
	stmt, err := this.db.Prepare("INSERT tags SET blogid=?,name=?")
	CheckErr(err)
	defer stmt.Close()

	_, err1 := stmt.Exec(blogid, name)
	CheckErr(err1)
}*/

func (this *Dao) DelTagByName(name []string) {
	stmt, err := this.db.Prepare("DELETE  FROM tags WHERE name =?")
	CheckErr(err)
	defer stmt.Close()

	finish := make(chan bool)
	lenth := len(name)
	for i := 0; i < lenth; i++ {
		go func(tagName string) {
			defer func() { finish <- true }()
			if _, err = stmt.Exec(tagName); err != nil {
				fmt.Println("stmt.Exec: ", err.Error())
				return
			}
		}(name[i])
	}

	for i := 0; i < lenth; i++ {
		<-finish
	}
}

func (this *Dao) DelTagByBlogId(id int) {
	stmt, err := this.db.Prepare("DELETE  FROM tags WHERE blogid =?")
	CheckErr(err)
	defer stmt.Close()

	_, err1 := stmt.Exec(id)
	CheckErr(err1)
}

func (this *Dao) FindByTag(page, limit int, name string) map[string]Blog {
	rows, err := this.db.Query("select id,title,content,created from blogs where id in ( select blogid from tags where name = ?) order by created desc limit ?,?", name, page, limit)
	CheckErr(err)
	defer rows.Close()
	blogs := make(map[string]Blog)
	for i := 0; rows.Next(); i++ {
		var blog = new(Blog)
		err = rows.Scan(&blog.Id, &blog.Title, &blog.Content, &blog.Created)
		CheckErr(err)
		key := fmt.Sprintf("%d", i)
		blogs[key] = *blog
	}

	return blogs
}

func (this *Dao) FindTag() map[string]Tag {
	rows, err := this.db.Query("select distinct name,id,blogid from tags group by name")
	CheckErr(err)
	defer rows.Close()
	tags := make(map[string]Tag)
	for i := 0; rows.Next(); i++ {
		var tag = new(Tag)
		err = rows.Scan(&tag.Name, &tag.Id, &tag.BlogId)
		CheckErr(err)
		key := fmt.Sprintf("%d", i)
		tags[key] = *tag
	}

	return tags
}
