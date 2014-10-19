package model

func (this *Dao) Click(blogid int) {
	stmt, err := this.db.Prepare("update ctrs SET count = count +1 where blogid = ? ")
	CheckErr(err)
	defer stmt.Close()

	_, err1 := stmt.Exec(blogid)
	CheckErr(err1)
}

func (this *Dao) AddCtr(blogid int) {
	stmt, err := this.db.Prepare("insert into ctrs SET blogid = ?  ,count = 0")
	CheckErr(err)
	defer stmt.Close()

	_, err1 := stmt.Exec(blogid)
	CheckErr(err1)
}
