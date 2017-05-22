package models

import (
	"time"
)

type Users struct {
	Id        int       `xorm:"not null pk autoincr INT(10)"`
	Name      string    `xorm:"not null default '' unique VARCHAR(32)"`
	Realname  string    `xorm:"not null default '' VARCHAR(32)"`
	Password  string    `xorm:"not null default '' CHAR(32)"`
	Salt      string    `xorm:"not null default '' CHAR(32)"`
	Sign      string    `xorm:"not null default '' CHAR(32)"`
	Email     string    `xorm:"not null default '' VARCHAR(255)"`
	GroupId   int       `xorm:"not null default 0 INT(10)"`
	IsDeleted string    `xorm:"not null default 'N' ENUM('Y','N')"`
	UpdatedAt time.Time `xorm:"not null default '0000-00-00 00:00:00' TIMESTAMP updated"`
	CreatedAt time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP created"`
}

func (this *Users) TableName() string {

	return "users"
}

func (this *Users) Save(t *Users) (err error) {

	if t.Id > 0 {
		err = dao.UpdateById(t)
	} else {
		err = dao.Insert(t)
	}

	return
}

func (this *Users) Update(params Cond,
	where Cond) (err error) {

	return dao.UpdateByWhere(this, params, where)
}

func (this *Users) SearchOne(cond Cond) (t *Users, err error) {

	cond["page"] = 1
	cond["pagesize"] = 1

	rs, err := this.Search(cond)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}

func (this *Users) Search(cond Cond) (t []*Users, err error) {

	err = dao.Search(&t, cond)

	return
}

func (this *Users) Count(cond Cond) (total int64, err error) {

	total, err = dao.Count(this, cond)

	return
}

func (this *Users) GetMulti(ids ...interface{}) (t []*Users, err error) {
	err = dao.GetMulti(&t, ids...)

	return
}

//注意，和SearchOne一样，返回的t可能是nil TODO
func (this *Users) GetById(id interface{}) (t *Users, err error) {

	rs, err := this.GetMulti(id)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	} else {
		return this, err
	}

	return
}
