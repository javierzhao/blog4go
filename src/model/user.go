package model

type User struct {
	Id        int
	User_name string
	Pass_word string
}

// 判断用户是否存在
func (this *Dao) IsUser(username, password string) bool {
	rows, err := this.db.Query("select id from users where user_name = ? and pass_word = ?", username, password)
	CheckErr(err)
	defer rows.Close()
	return rows.Next()
}

/**
func (this *Dao) Insert(name, pass string) {
	//插入数据
	stmt, err := this.db.Prepare("INSERT users SET user_name=?,pass_word=?")
	checkErr(err)
	defer stmt.Close()

	res, err := stmt.Exec(name, pass)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	// return id
	fmt.Println(id)

	// db.Close()
}*/
