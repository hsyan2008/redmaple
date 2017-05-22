package models

import (
	"time"
)

type Messages struct {
	Id        int       `xorm:"not null pk autoincr INT(10)"`
	TaskId    int       `xorm:"not null default 0 INT(10)"`
	UserId    int       `xorm:"not null default 0 INT(10)"`
	Status    int       `xorm:"not null default 0 TINYINT(3)"`
	Msg       string    `xorm:"not null default '' VARCHAR(255)"`
	IsDeleted string    `xorm:"not null default 'N' ENUM('Y','N')"`
	UpdatedAt time.Time `xorm:"not null default '0000-00-00 00:00:00' TIMESTAMP updated"`
	CreatedAt time.Time `xorm:"not null default 'CURRENT_TIMESTAMP' TIMESTAMP created"`
	Task      *Tasks    `xorm:"-"`
	User      *Users    `xorm:"-"`
}

func (this *Messages) TableName() string {

	return "messages"
}

func (this *Messages) Save(t *Messages) (err error) {

	if t.Id > 0 {
		err = dao.UpdateById(t)
	} else {
		err = dao.Insert(t)
	}

	return
}

func (this *Messages) Update(params Cond,
	where Cond) (err error) {

	return dao.UpdateByWhere(this, params, where)
}

func (this *Messages) SearchOne(cond Cond) (t *Messages, err error) {

	cond["page"] = 1
	cond["pagesize"] = 1

	rs, err := this.Search(cond)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	}

	return
}

func (this *Messages) Search(cond Cond) (t []*Messages, err error) {

	err = dao.Search(&t, cond)

	for k, v := range t {
		t[k].Task, _ = taskModel.GetById(v.TaskId)
		t[k].User, _ = userModel.GetById(v.UserId)
	}

	return
}

func (this *Messages) Count(cond Cond) (total int64, err error) {

	total, err = dao.Count(this, cond)

	return
}

func (this *Messages) GetMulti(ids ...interface{}) (t []*Messages, err error) {
	err = dao.GetMulti(&t, ids...)

	for k, v := range t {
		t[k].Task, _ = taskModel.GetById(v.TaskId)
		t[k].User, _ = userModel.GetById(v.UserId)
	}

	return
}

//注意，和SearchOne一样，返回的t可能是nil TODO
func (this *Messages) GetById(id interface{}) (t *Messages, err error) {

	rs, err := this.GetMulti(id)
	if err == nil && len(rs) > 0 {
		t = rs[0]
	} else {
		return this, err
	}

	return
}
