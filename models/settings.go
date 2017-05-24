package models

import (
	"time"
)

type Settings struct {
	Id        int       `xorm:"not null pk autoincr INT(10)"`
	SmtpAddr  string    `xorm:"not null default '' VARCHAR(128)"`
	SmtpPort  int       `xorm:"not null default 0 SMALLINT(5)"`
	SmtpSsl   string    `xorm:"not null default 'N' ENUM('Y','N')"`
	SmtpName  string    `xorm:"not null default '' VARCHAR(32)"`
	SmtpUser  string    `xorm:"not null default '' VARCHAR(128)"`
	SmtpPass  string    `xorm:"not null default '' VARCHAR(64)"`
	IsDeleted string    `xorm:"not null default 'N' ENUM('Y','N')"`
	UpdatedAt time.Time `xorm:"not null default '0000-00-00 00:00:00' TIMESTAMP updated"`
	CreatedAt time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP created"`
}

func (this *Settings) TableName() string {

	return "settings"
}

func (this *Settings) Save(t *Settings) (err error) {

	if t.Id > 0 {
		err = dao.UpdateById(t)
	} else {
		err = dao.Insert(t)
	}

	return
}

func (this *Settings) Update(params Cond,
	where Cond) (err error) {

	return dao.UpdateByWhere(this, params, where)
}

func (this *Settings) SearchOne(cond Cond) (t *Settings, err error) {

	cond["page"] = 1
	cond["pagesize"] = 1

	rs, err := this.Search(cond)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}

func (this *Settings) Search(cond Cond) (t []*Settings, err error) {

	err = dao.Search(&t, cond)

	return
}

func (this *Settings) Count(cond Cond) (total int64, err error) {

	total, err = dao.Count(this, cond)

	return
}

func (this *Settings) GetMulti(ids ...interface{}) (t []*Settings, err error) {
	err = dao.GetMulti(&t, ids...)

	return
}

//注意，和SearchOne一样，返回的t可能是nil TODO
func (this *Settings) GetById(id interface{}) (t *Settings, err error) {

	rs, err := this.GetMulti(id)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	} else {
		return this, err
	}

	return
}
