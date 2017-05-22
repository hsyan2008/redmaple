package models

import (
	"time"
)

type Machines struct {
	Id        int       `xorm:"not null pk autoincr INT(10)"`
	Name      string    `xorm:"not null default '' VARCHAR(32)"`
	Env       int       `xorm:"not null default 0 TINYINT(3)"`
	Ip        string    `xorm:"not null default '' VARCHAR(16)"`
	Port      string    `xorm:"not null default '' VARCHAR(5)"`
	User      string    `xorm:"not null default '' VARCHAR(32)"`
	Auth      string    `xorm:"not null TEXT"`
	InnerIp   string    `xorm:"not null default '' VARCHAR(16)"`
	InnerPort string    `xorm:"not null default '' VARCHAR(5)"`
	InnerUser string    `xorm:"not null default '' VARCHAR(32)"`
	InnerAuth string    `xorm:"not null TEXT"`
	IsDeleted string    `xorm:"not null default 'N' ENUM('Y','N')"`
	UpdatedAt time.Time `xorm:"not null default '0000-00-00 00:00:00' TIMESTAMP updated"`
	CreatedAt time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP created"`
}

func (this *Machines) TableName() string {

	return "machines"
}

func (this *Machines) Save(t *Machines) (err error) {

	if t.Id > 0 {
		err = dao.UpdateById(t)
	} else {
		err = dao.Insert(t)
	}

	return
}

func (this *Machines) Update(params Cond,
	where Cond) (err error) {

	return dao.UpdateByWhere(this, params, where)
}

func (this *Machines) SearchOne(cond Cond) (t *Machines, err error) {

	cond["page"] = 1
	cond["pagesize"] = 1

	rs, err := this.Search(cond)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}

func (this *Machines) Search(cond Cond) (t []*Machines, err error) {

	err = dao.Search(&t, cond)

	return
}

func (this *Machines) Count(cond Cond) (total int64, err error) {

	total, err = dao.Count(this, cond)

	return
}

func (this *Machines) GetMulti(ids ...interface{}) (t []*Machines, err error) {
	err = dao.GetMulti(&t, ids...)

	return
}

//注意，和SearchOne一样，返回的t可能是nil TODO
func (this *Machines) GetById(id interface{}) (t *Machines, err error) {

	rs, err := this.GetMulti(id)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}
