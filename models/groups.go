package models

import (
	"time"
)

type Groups struct {
	Id        int       `xorm:"not null pk autoincr INT(10)"`
	Name      string    `xorm:"not null unique VARCHAR(32)"`
	IsDeleted string    `xorm:"not null default 'N' ENUM('Y','N')"`
	UpdatedAt time.Time `xorm:"not null default '0000-00-00 00:00:00' TIMESTAMP updated"`
	CreatedAt time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP"`
}

func (this *Groups) TableName() string {

	return "groups"
}

func (this *Groups) Save(t *Groups) (err error) {

	if t.Id > 0 {
		err = dao.UpdateById(t)
	} else {
		err = dao.Insert(t)
	}

	return
}

func (this *Groups) Update(params Cond,
	where Cond) (err error) {

	return dao.UpdateByWhere(this, params, where)
}

func (this *Groups) SearchOne(cond Cond) (t *Groups, err error) {

	cond["page"] = 1
	cond["pagesize"] = 1

	rs, err := this.Search(cond)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}

func (this *Groups) Search(cond Cond) (t []*Groups, err error) {

	err = dao.Search(&t, cond)

	return
}

func (this *Groups) Count(cond Cond) (total int64, err error) {

	total, err = dao.Count(this, cond)

	return
}

func (this *Groups) GetMulti(ids ...interface{}) (t []*Groups, err error) {
	err = dao.GetMulti(&t, ids...)

	return
}

//注意，和SearchOne一样，返回的t可能是nil TODO
func (this *Groups) GetById(id interface{}) (t *Groups, err error) {

	rs, err := this.GetMulti(id)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}
